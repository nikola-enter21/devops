package middleware

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/nikola-enter21/devops-fmi-course/internal/logging"
	"github.com/nikola-enter21/devops-fmi-course/internal/policy"
)

var (
	log = logging.MustNewLogger()
)

type Authorizer interface {
	Authorize(ctx context.Context, role, route string) (bool, error)
}

type OPAAuthorizer struct {
	engine *policy.Engine
}

func NewOPAAuthorizer(engine *policy.Engine) Authorizer {
	return &OPAAuthorizer{
		engine: engine,
	}
}

func (a *OPAAuthorizer) Authorize(ctx context.Context, role, route string) (bool, error) {
	return a.engine.Evaluate(ctx, role, route)
}

func AuthorizeMiddleware(auth Authorizer) fiber.Handler {
	return func(c fiber.Ctx) error {
		role := c.Get("X-User-Role") // demo purpose only
		route := c.Route().Name
		ctx := c.Context()

		if route == "" {
			log.Warnw("unauthorized route", "path", c.Path())
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "route has no name, cannot authorize",
			})
		}

		allowed, err := auth.Authorize(ctx, role, route)
		if err != nil {
			log.Errorw("OPA evaluation error", "route", route, "role", role, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "authorization engine failure",
			})
		}

		if !allowed {
			log.Warnw("access denied", "role", role, "route", route, "path", c.Path(), "method", c.Method())
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": fmt.Sprintf("access denied for role %s on %s", role, route),
			})
		}

		log.Infow("access granted", "role", role, "route", route, "path", c.Path(), "method", c.Method())
		return c.Next()
	}
}
