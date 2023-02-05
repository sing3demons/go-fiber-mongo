package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/sing3demons/go-fiber-mongo/helper"
	"github.com/sing3demons/go-fiber-mongo/models"
	"github.com/sing3demons/go-fiber-mongo/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type categoryResponse struct {
// 	Name    string `json:"name" bson:"_id,omitempty"`
// 	Product []struct {
// 		ID   primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
// 		Name string             `json:"name" bson:"name,omitempty"`
// 	} `json:"products"`
// }

type allCategoryResponse struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name,omitempty"`
}

type createCategoryForm struct {
	Name string `json:"name" bson:"name,omitempty" validate:"required"`
}

type CategoryController interface {
	FindCategories(ctx *fiber.Ctx) error
	CreateCategory(c *fiber.Ctx) error
}

type categoryController struct {
	Service service.CategoryService
}

func NewCategoryController(service service.CategoryService) CategoryController {
	return &categoryController{Service: service}
}

func (tx *categoryController) FindCategories(c *fiber.Ctx) error {
	categories, err := tx.Service.FindAll()
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Something went wrong"})
	}

	serializedCategory := []allCategoryResponse{}
	copier.Copy(&serializedCategory, &categories)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"category": serializedCategory})
}

func (tx *categoryController) CreateCategory(c *fiber.Ctx) error {
	var form createCategoryForm
	if err := c.BodyParser(&form); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	if err := helper.ValidateStruct(&form); err != nil {
		return c.JSON(err)
	}

	var category models.Category

	copier.Copy(&category, &form)
	result, err := tx.Service.Create(category)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Something went wrong"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"category": result})
}
