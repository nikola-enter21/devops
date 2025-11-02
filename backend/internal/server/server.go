package server

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"github.com/nikola-enter21/devops-fmi-course/internal/handlers"
	"github.com/nikola-enter21/devops-fmi-course/internal/middleware"
)

type Server struct {
	App        *fiber.App
	Authorizer middleware.Authorizer
}

func NewServer(authorizer middleware.Authorizer) *Server {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))
	app.Use(logger.New())
	app.Use(recover.New())

	s := &Server{
		App:        app,
		Authorizer: authorizer,
	}

	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	s.App.Route("/api/v1", func(api fiber.Router) {
		wrapRoutes(api, []fiber.Handler{middleware.AuthorizeMiddleware(s.Authorizer)}, func(r fiber.Router) {
			r.Get("/healthz", handlers.HealthCheckHandler).Name("healthcheck")
			r.Get("/checkDatabase", handlers.CheckDatabaseHandler(
				getEnv("DB_HOST", "postgres"),
				getEnv("DB_PORT", "5432"),
				getEnv("DB_USER", "postgres"),
				getEnv("DB_PASSWORD", "postgres"),
				getEnv("DB_NAME", "postgres"),
			)).Name("checkDatabase")

			r.Post("/login", handlers.LoginHandler).Name("auth.login")
			r.Post("/register", handlers.RegisterHandler).Name("auth.register")
		})
	}, "api.v1.")
}

func (s *Server) Listen(port string) error {
	return s.App.Listen(":" + port)
}
