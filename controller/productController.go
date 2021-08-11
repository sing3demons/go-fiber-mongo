package controller

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/sing3demons/go-fiber-mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type createProduct struct {
	Name  string `form:"name" bson:"name" validate:"required"`
	Desc  string `form:"desc" bson:"desc" validate:"required"`
	Price int    `form:"price" bson:"price" validate:"required"`
	Image string `form:"image" bson:"image" validate:"required"`
}

type updateProduct struct {
	Name  string `form:"name" bson:"name"`
	Desc  string `form:"desc" bson:"desc"`
	Price int    `form:"price" bson:"price"`
	Image string `form:"image" bson:"image"`
}

type productRespons struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name"`
	Desc  string             `json:"desc"`
	Price int                `json:"price"`
	Image string             `json:"image"`
}

type productController struct {
	db *mongo.Database
}

type ProductController interface {
	FindProducts(c *fiber.Ctx) error
	FindProduct(c *fiber.Ctx) error
	CreateProduct(c *fiber.Ctx) error
	UpdataProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
}

func NewItemController(db *mongo.Database) ProductController {
	return &productController{db: db}
}

func (tx *productController) collection() *mongo.Collection {
	return tx.db.Collection(("products"))
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter, err := tx.findProductByID(c)
	if err != nil {
		c.Status(fiber.StatusNotFound).JSON(err)
	}
	var product models.Product

	if err := tx.collection().FindOne(ctx, filter).Decode(&product); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"product": product})
}

func (tx *productController) FindProducts(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := tx.collection().Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}
	defer cursor.Close(ctx)
	products := []models.Product{}

	if err := cursor.All(ctx, &products); err != nil {
		panic(err)
	}

	serializedProduct := []productRespons{}
	copier.Copy(&serializedProduct, &products)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"products": serializedProduct})
}

func (tx *productController) CreateProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var form createProduct
	if err := c.BodyParser(&form); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	var product models.Product

	copier.Copy(&product, &form)

	image, err := tx.setProductImage(c, &product)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}

	product.Image = image

	_, err = tx.collection().InsertOne(ctx, product)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}

	// fmt.Println(result.InsertedID)
	serializedProduct := productRespons{}
	copier.Copy(&serializedProduct, &product)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"product": serializedProduct})
}

func (tx *productController) UpdataProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var form updateProduct
	if err := c.BodyParser(&form); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	filter, err := tx.findProductByID(c)
	if err != nil {
		c.Status(fiber.StatusNotFound).JSON(err)
	}

	var product models.Product

	if err := tx.collection().FindOne(ctx, filter).Decode(&product); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
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

	image, _ := tx.setProductImage(c, &product)
	form.Image = image

	if form.Image == "" {
		form.Image = product.Image
	}

	// copier.Copy(&product, &form)

	update := bson.D{
		{"$set", form},
	}

	if err := tx.collection().FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		return c.JSON(err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (tx *productController) DeleteProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter, err := tx.findProductByID(c)
	if err != nil {
		c.Status(fiber.StatusNotFound).JSON(err)
	}

	if err := tx.removeImageProduct(filter); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}

	if err := tx.collection().FindOneAndDelete(ctx, filter).Err(); err != nil {
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

	filename := "uploads/products" + "/" + "images" + "/" + strings.Replace(uniqueId.String(), "-", "", -1)
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var product models.Product

	if err := tx.collection().FindOne(ctx, filter).Decode(&product); err != nil {
		return err
	}

	if product.Image != "" {
		image := product.Image
		fmt.Println(image)
		pwd, _ := os.Getwd()
		fmt.Println(pwd)
		os.Remove(pwd + "/" + image)
	}
	return nil
}
