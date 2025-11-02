package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/nikola-enter21/devops-fmi-course/logging"
	"github.com/nikola-enter21/devops-fmi-course/policy"
)

var (
	opaEngine *policy.Engine
	log       = logging.MustNewLogger()
)

func init() {
	engine, err := policy.NewEmbedded()
	if err != nil {
		log.Fatalf("failed to initialize OPA policy engine: %v", err)
	}
	opaEngine = engine
}

func AuthorizeMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		role := c.Get("X-User-Role") // For demo; normally from session/JWT
		route := c.Route().Name

		if route == "" {
			log.Warnw("unauthorized route", "path", c.Path())
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "route has no name, cannot authorize",
			})
		}

		// Evaluate OPA policy
		allowed, err := opaEngine.Evaluate(role, route)
		if err != nil {
			log.Errorw("OPA evaluation error",
				"route", route,
				"role", role,
				"error", err,
			)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "authorization engine failure",
			})
		}

		if !allowed {
			log.Warnw("access denied",
				"role", role,
				"route", route,
				"path", c.Path(),
				"method", c.Method(),
			)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": fmt.Sprintf("access denied for role %s on %s", role, route),
			})
		}

		log.Infow("access granted",
			"role", role,
			"route", route,
			"path", c.Path(),
			"method", c.Method(),
		)

		return c.Next()
	}
}
