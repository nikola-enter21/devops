package handlers

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v3"
	_ "github.com/lib/pq"
	"github.com/nikola-enter21/devops-fmi-course/internal/logging"
)

var (
	log = logging.MustNewLogger()
)

func HealthCheckHandler(c fiber.Ctx) error {
	log.Infow("health check", "ip", c.IP())
	return c.SendStatus(fiber.StatusOK)
}

func LoginHandler(c fiber.Ctx) error {
	log.Infow("login attempt", "ip", c.IP(), "userAgent", string(c.Request().Header.UserAgent()))
	return c.JSON(fiber.Map{
		"message": "Login successful",
	})
}

func RegisterHandler(c fiber.Ctx) error {
	log.Infow("user registration", "ip", c.IP())
	return c.JSON(fiber.Map{
		"message": "User registered",
	})
}

func CheckDatabaseHandler(host, port, user, password, dbname string) fiber.Handler {
	return func(c fiber.Ctx) error {
		connStr := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname,
		)

		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Errorw("failed to open DB connection", "error", err)
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer db.Close()

		if err := db.Ping(); err != nil {
			log.Errorw("database unreachable", "error", err)
			return c.Status(500).JSON(fiber.Map{
				"status": "unreachable",
				"error":  err.Error(),
			})
		}

		log.Infow("database connection successful",
			"host", host,
			"port", port,
			"user", user,
			"dbname", dbname,
		)

		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Connected to database",
		})
	}
}
