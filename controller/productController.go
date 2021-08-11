package controller

import (
	"context"
	"fmt"
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
	FindProduct(c *fiber.Ctx) error
	CreateProduct(c *fiber.Ctx) error
}

func NewItemController(db *mongo.Database) ProductController {
	return &productController{db: db}
}

func (tx *productController) collection() *mongo.Collection {
	return tx.db.Collection(("products"))
}

func (tx *productController) FindProduct(c *fiber.Ctx) error {
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

func (tx *productController) setProductImage(c *fiber.Ctx, product *models.Product) (string, error) {

	file, err := c.FormFile("image")
	if err != nil || file == nil {
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
