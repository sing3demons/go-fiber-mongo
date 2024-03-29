package controller

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/sing3demons/go-fiber-mongo/models"
	"github.com/sing3demons/go-fiber-mongo/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createProduct struct {
	Name       string `form:"name" bson:"name" validate:"required"`
	Desc       string `form:"desc" bson:"desc" validate:"required"`
	Price      int    `form:"price" bson:"price" validate:"required"`
	Image      string `form:"image" bson:"image" validate:"required"`
	CategoryID string `form:"categoryId" bson:"categoryId,omitempty"`
}

type updateProduct struct {
	Name       string `form:"name" bson:"name"`
	Desc       string `form:"desc" bson:"desc"`
	Price      int    `form:"price" bson:"price"`
	Image      string `form:"image" bson:"image"`
	CategoryID string `form:"categoryId,omitempty" bson:"categoryId,omitempty"`
}

type productResponse struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name"`
	Desc       string             `json:"desc"`
	Price      int                `json:"price"`
	Image      string             `json:"image"`
	CategoryID primitive.ObjectID `json:"categoryId" bson:"categoryId,omitempty"`
	Category   struct {
		ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		Name string             `json:"name"`
	} `json:"category"`
}

type productController struct {
	Service service.ProductService
}

type ProductController interface {
	FindProducts(c *fiber.Ctx) error
	FindProduct(c *fiber.Ctx) error
	CreateProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
}

func NewItemController(service service.ProductService) ProductController {
	return &productController{Service: service}
}

func (tx *productController) findProductByID(c *fiber.Ctx) (primitive.M, error) {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	return filter, nil
}

func (tx *productController) FindProduct(c *fiber.Ctx) error {

	filter, err := tx.findProductByID(c)
	if err != nil {
		c.Status(fiber.StatusNotFound).JSON(err)
	}

	product, err := tx.Service.FindOne(filter)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "recods not found"})
	}

	serializedProduct := productResponse{}
	copier.Copy(&serializedProduct, &product)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"product": serializedProduct})
}

func (tx *productController) FindProducts(c *fiber.Ctx) error {
	products, err := tx.Service.FindAll()

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Something went wrong"})
	}

	serializedProduct := []productResponse{}
	copier.Copy(&serializedProduct, &products)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"products": serializedProduct})
}

// CreateProduct
func (tx *productController) CreateProduct(c *fiber.Ctx) error {
	var form createProduct
	if err := c.BodyParser(&form); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	var product models.Product

	copier.Copy(&product, &form)
	product.CategoryID, _ = primitive.ObjectIDFromHex(form.CategoryID)

	image, err := tx.setProductImage(c, &product)
	if err != nil {

		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}

	product.Image = image

	result, err := tx.Service.Create(product)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Something went wrong"})
	}

	// fmt.Println(result.InsertedID)
	serializedProduct := productResponse{}
	copier.Copy(&serializedProduct, &result)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"product": serializedProduct})
}

func (tx *productController) UpdateProduct(c *fiber.Ctx) error {
	var form updateProduct
	if err := c.BodyParser(&form); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	filter, err := tx.findProductByID(c)
	if err != nil {
		c.Status(fiber.StatusNotFound).JSON(err)
	}

	product, err := tx.Service.FindOne(filter)
	if err != nil {
		c.Status(fiber.StatusNotFound).JSON(err)
	}

	//validate
	if form.Name == "" {
		form.Name = product.Name
	}

	if form.Desc == "" {
		form.Desc = product.Desc
	}

	if form.Price == 0 {
		form.Price = product.Price
	}

	image, _ := tx.setProductImage(c, product)
	form.Image = image

	if form.Image == "" {
		form.Image = product.Image
	}

	copier.Copy(&product, &form)
	product.CategoryID, _ = primitive.ObjectIDFromHex(form.CategoryID)

	// update := bson.D{
	// 	{"$set", form},
	// }

	update := []interface{}{bson.D{
		{Key: "$set", Value: form},
	}}

	if err := tx.removeImageProduct(filter); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}

	if err := tx.Service.Update(filter, update); err != nil {
		return c.JSON(err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (tx *productController) DeleteProduct(c *fiber.Ctx) error {

	filter, err := tx.findProductByID(c)
	if err != nil {
		c.Status(fiber.StatusNotFound).JSON(err)
	}

	if err := tx.removeImageProduct(filter); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}

	if err := tx.Service.Delete(filter); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (tx *productController) setProductImage(c *fiber.Ctx, product *models.Product) (string, error) {
	file, err := c.FormFile("image")
	if err != nil || file == nil {
		log.Println(err)
		return "", err
	}

	// generate new uuid for image name
	uniqueId := uuid.New()
	filename := "/uploads/products" + "/" + "images" + "/" + strings.Replace(uniqueId.String(), "-", "", -1)
	// extract image extension from original file filename
	fileExt := strings.Split(file.Filename, ".")[1]
	// generate image from filename and extension
	image := fmt.Sprintf("%s.%s", filename, fileExt)

	if err := c.SaveFile(file, image); err != nil {
		return "", err
	}

	return image, nil
}

func (tx *productController) removeImageProduct(filter primitive.M) error {

	product, err := tx.Service.FindOne(filter)

	if err != nil {
		return err
	}

	if product.Image != "" {
		image := product.Image
		pwd, _ := os.Getwd()
		os.Remove(pwd + "/" + image)
	}
	return nil
}
