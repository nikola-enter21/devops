package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikola-enter21/devops-fmi-course/authorizer"
	"github.com/nikola-enter21/devops-fmi-course/config"
	"github.com/nikola-enter21/devops-fmi-course/logging"
	"github.com/nikola-enter21/devops-fmi-course/service/db/repo"
	"github.com/nikola-enter21/devops-fmi-course/service/grpc"
)

var (
	log      = logging.MustNewLogger()
	httpPort = config.Env("HTTP_PORT", "8080")
	grpcPort = config.Env("GRPC_PORT", "8079")
)

func main() {
	defer logging.Sync()

	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	auth, err := authorizer.NewEmbedded()
	if err != nil {
		log.Fatalf("failed to initialize OPA authorizer: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), config.PostgresDSN())
	if err != nil {
		log.Fatalf("failed to create a pool: %v", err)
	}
	defer pool.Close()

	s := &grpc.Server{
		Authorizer:     auth,
		UserRepository: repo.NewUserRepository(pool),
	}

	s.Serve(signalCtx, httpPort, grpcPort)
}
