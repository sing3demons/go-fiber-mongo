package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/go-fiber-mongo/controller"
	"github.com/sing3demons/go-fiber-mongo/database"
	"github.com/sing3demons/go-fiber-mongo/repository"
	"github.com/sing3demons/go-fiber-mongo/service"
)

func Serve(app *fiber.App) {
	db := database.InitDB()
	cache := database.NewRedisCache("redis:6379", 1, 10)
	v1 := app.Group("api/v1")

	productGroup := v1.Group("/products")
	productRepository := repository.NewProductRepository(db, cache)
	productSevice := service.NewProductService(productRepository)
	productController := controller.NewItemController(productSevice)
	{
		productGroup.Get("", productController.FindProducts)
		productGroup.Get("/:id", productController.FindProduct)
		productGroup.Delete("/:id", productController.DeleteProduct)
		productGroup.Put("/:id", productController.UpdateProduct)
		productGroup.Post("", productController.CreateProduct)
	}
}
