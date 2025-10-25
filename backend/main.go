package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("ALLOWED_ORIGINS"),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Use(logger.New())
	app.Use(recover.New())

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Login successful",
		})
	})

	app.Post("/register", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "User registered",
		})
	})

	port := getEnv("PORT", "8080")
	log.Printf("Server running on http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
