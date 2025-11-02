package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"github.com/nikola-enter21/devops-fmi-course/handlers"
	"github.com/nikola-enter21/devops-fmi-course/logging"
	"github.com/nikola-enter21/devops-fmi-course/middleware"
)

func main() {
	log := logging.MustNewLogger()
	defer log.Sync()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	app.Use(logger.New())
	app.Use(recover.New())

	app.Route("/api/v1", func(api fiber.Router) {
		wrapRoutes(api, []fiber.Handler{middleware.AuthorizeMiddleware()}, func(r fiber.Router) {
			r.Get("/healthz", handlers.HealthCheckHandler).Name("healthcheck")
			r.Get("/checkDatabase", handlers.CheckDatabaseHandler(
				getEnv("DB_HOST", "postgres"),
				getEnv("DB_PORT", "5432"),
				getEnv("DB_USER", "postgres"),
				getEnv("DB_PASSWORD", "postgres"),
				getEnv("DB_NAME", "postgres"),
			)).Name("healthcheck")

			r.Post("/login", handlers.LoginHandler).Name("auth.login")
			r.Post("/register", handlers.RegisterHandler).Name("auth.register")
		})
	}, "api.v1.")

	port := getEnv("PORT", "8080")

	log.Infow("server starting", "url", "http://localhost:"+port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalw("server failed to start", "error", err)
	}
}
