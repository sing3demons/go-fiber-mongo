package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/sing3demons/go-fiber-mongo/routes"
)

var (
	buildcommit = "dev"
	buildtime   = time.Now().String()
)

func main() {
	// Liveness Probe
	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("/tmp/live")

	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	app := fiber.New()
	app.Use(recover.New())
	app.Static("/", "./wwwroot")
	// Readiness Probe
	app.Get("/healthz", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	app.Get("/dashboard", monitor.New())
	app.Get("/x", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"buildcommit": buildcommit,
			"buildtime":   buildtime,
		})
	})

	app.Static("/uploads", "./uploads")

	//สร้าง folder
	uploadDirs := [...]string{"products", "users"}

	for _, dir := range uploadDirs {
		path := fmt.Sprintf("uploads/%s/images", dir)
		os.MkdirAll(path, 0755)
	}

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hello, World 👋!")
	// })

	routes.Serve(app)

	//Graceful Shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := app.Listen(":" + os.Getenv("PORT")); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()

	fmt.Println("shutting down gracefully, press Ctrl+C again to force")

	if err := app.Shutdown(); err != nil {
		fmt.Println(err)
	}

}
