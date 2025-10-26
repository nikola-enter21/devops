package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"
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

	app.Get("/checkDatabase", func(c *fiber.Ctx) error {
		connStr := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			getEnv("DB_HOST", "postgres"),
			getEnv("DB_PORT", "5432"),
			getEnv("DB_USER", "postgres"),
			getEnv("DB_PASSWORD", "postgres"),
			getEnv("DB_NAME", "postgres"),
		)

		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Println("DB connection open error:", err)
			return c.Status(501).JSON(fiber.Map{"error": err.Error()})
		}
		defer db.Close()

		if err := db.Ping(); err != nil {
			log.Println("DB ping error:", err)
			return c.Status(500).JSON(fiber.Map{
				"status": "unreachable",
				"error":  err.Error(),
			})
		}

		log.Println("Connected to Postgres successfully")
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Connected to database",
		})
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
