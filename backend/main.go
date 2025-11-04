package main

import (
	"context"
	"fmt"
	"net"
	"os/signal"
	"syscall"

	"buf.build/go/protovalidate"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	user "github.com/nikola-enter21/devops-fmi-course/api/gen/go/user/v1"
	"github.com/nikola-enter21/devops-fmi-course/authorizer"
	"github.com/nikola-enter21/devops-fmi-course/gateway"
	"github.com/nikola-enter21/devops-fmi-course/logging"
	"github.com/nikola-enter21/devops-fmi-course/service"
	"google.golang.org/grpc"
)

var (
	log      = logging.MustNewLogger()
	httpPort = getEnv("HTTP_PORT", "8080")
	grpcPort = getEnv("GRPC_PORT", "8079")
)

type Server struct {
	Authorizer authorizer.Authorizer
}

func main() {
	auth, err := authorizer.NewEmbedded()
	if err != nil {
		log.Fatalf("failed to initialize OPA authorizer: %v", err)
	}

	s := &Server{Authorizer: auth}
	s.Serve()
}

func (s *Server) Serve() {
	defer logging.Sync()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	validator, err := protovalidate.New()
	if err != nil {
		log.Fatalf("failed to create proto validator: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			protovalidate_middleware.UnaryServerInterceptor(validator),
			authorizer.UnaryServerInterceptor(s.Authorizer),
		),
		grpc.ChainStreamInterceptor(
			protovalidate_middleware.StreamServerInterceptor(validator),
			authorizer.StreamServerInterceptor(s.Authorizer),
		),
	)

	userSvc := &service.UserServiceServer{}
	user.RegisterUserServiceServer(grpcServer, userSvc)

	// Start the HTTP gateway in a separate goroutine
	gatewayDone := make(chan struct{})
	go func() {
		gateway.Serve(ctx, ":"+httpPort, fmt.Sprintf("localhost%s", ":"+grpcPort))
		close(gatewayDone)
	}()

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen", "error", err)
	}
	log.Infow("gRPC server listening", "port", ":"+grpcPort)

	// Watch for shutdown signal and gracefully stop the gRPC server.
	go func() {
		<-ctx.Done()
		log.Infow("gRPC shutdown signal received, waiting for gateway to stop...")

		// Wait for the gateway to stop first.
		<-gatewayDone

		log.Infow("Stopping gRPC gracefully...")
		grpcServer.GracefulStop()
	}()

	// Blocks here until the gRPC server is completely stopped.
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server exited", "error", err)
	}

	log.Infow("gRPC server stopped cleanly")
}
