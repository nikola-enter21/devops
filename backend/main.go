package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/nikola-enter21/devops/db"
	"golang.org/x/crypto/bcrypt"
)

func normalizeEmail(e string) string {
	return strings.TrimSpace(strings.ToLower(e))
}

func readCreds(c *fiber.Ctx) (email, password string) {
	email = c.FormValue("email")
	password = c.FormValue("password")
	if email != "" || password != "" {
		return
	}
	var m map[string]string
	_ = json.Unmarshal(c.Body(), &m)
	email = m["email"]
	password = m["password"]
	return
}

func main() {

	app := fiber.New()
	app.Use(logger.New())
	mem := db.NewStore()

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"ok": true})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/users", func(c *fiber.Ctx) error {
		email, pw := readCreds(c)
		email = normalizeEmail(email)
		if email == "" || len(pw) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email and password required"})
		}
		if len(pw) < 8 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "password too short (min 8)"})
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(pw), 12)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "hash error"})
		}
		if !mem.Add(email, hash) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already exists"})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"email": email})
	})

	app.Post("/auth/login", func(c *fiber.Ctx) error {
		email, pw := readCreds(c)
		email = normalizeEmail(email)
		hash, ok := mem.Get(email)
		if !ok || bcrypt.CompareHashAndPassword(hash, []byte(pw)) != nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"ok": false})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"ok": true})
	})

	port := os.Getenv("PORT")
	log.Fatal(app.Listen(":" + port))
}
