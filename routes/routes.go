package routes

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/go-fiber-mongo/controller"
	"github.com/sing3demons/go-fiber-mongo/database"
	"github.com/sing3demons/go-fiber-mongo/repository"
	"github.com/sing3demons/go-fiber-mongo/service"
)

func Serve(app *fiber.App) {
	db := database.InitDB()
	host:=fmt.Sprintf("%s:6379",os.Getenv("REDIS_HOST"))
	cache := database.NewRedisCache(host, 0, 10)
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

	categoryGroup := v1.Group("/categories")
	categoryRepository := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryController := controller.NewCategoryController(categoryService)
	{
		categoryGroup.Get("", categoryController.FindCategories)
		categoryGroup.Post("", categoryController.CreateCategory)
	}
}
