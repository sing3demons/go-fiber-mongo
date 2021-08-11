package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/go-fiber-mongo/controller"
	"github.com/sing3demons/go-fiber-mongo/database"
)

func Serve(app *fiber.App) {
	db := database.InitDB()
	v1 := app.Group("api/v1")
	productController := controller.NewItemController(db)

	productGroup := v1.Group("/products")
	{
		productGroup.Get("", productController.FindProduct)
		productGroup.Post("", productController.CreateProduct)
	}
}