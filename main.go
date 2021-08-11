package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sing3demons/go-fiber-mongo/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := fiber.New()

	app.Static("/uploads", "./uploads")

	//สร้าง folder
	uploadDirs := [...]string{"products", "users"}

	for _, dir := range uploadDirs {
		path := fmt.Sprintf("uploads/%s/images", dir)
		os.MkdirAll(path, 0755)
	}
	
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	routes.Serve(app)
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
