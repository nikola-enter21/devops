package grpc

import (
	"context"
	"fmt"
	"net"

	"buf.build/go/protovalidate"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	_ "github.com/lib/pq"
	"github.com/nikola-enter21/devops-fmi-course/api/gen/go/user/v1"
	"github.com/nikola-enter21/devops-fmi-course/authorizer"
	"github.com/nikola-enter21/devops-fmi-course/gateway"
	"github.com/nikola-enter21/devops-fmi-course/logging"
	db "github.com/nikola-enter21/devops-fmi-course/service/db/gen"
	"google.golang.org/grpc"
)

var (
	log = logging.MustNewLogger()
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (db.User, error)
}

type Server struct {
	user.UnimplementedUserServiceServer

	Authorizer     authorizer.Authorizer
	UserRepository UserRepository
}

func (s *Server) Serve(ctx context.Context, httpPort, grpcPort string) {
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

	user.RegisterUserServiceServer(grpcServer, s)

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
