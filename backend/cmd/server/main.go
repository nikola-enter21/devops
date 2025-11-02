package main

import (
	"os"

	"github.com/nikola-enter21/devops-fmi-course/internal/logging"
	"github.com/nikola-enter21/devops-fmi-course/internal/middleware"
	"github.com/nikola-enter21/devops-fmi-course/internal/policy"
	"github.com/nikola-enter21/devops-fmi-course/internal/server"
)

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	log := logging.MustNewLogger()
	defer log.Sync()

	engine, err := policy.NewEmbedded()
	if err != nil {
		log.Fatalw("failed to initialize OPA engine", "error", err)
	}

	authorizer := middleware.NewOPAAuthorizer(engine)
	srv := server.NewServer(authorizer)

	port := getEnv("PORT", "8080")
	log.Infow("server starting", "url", "http://localhost:"+port)

	if err := srv.Listen(port); err != nil {
		log.Fatalw("server failed to start", "error", err)
	}
}
