package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"buf.build/go/protovalidate"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	_ "github.com/lib/pq"
	"github.com/nikola-enter21/devops-fmi-course/api/gen/go/user/v1"
	"github.com/nikola-enter21/devops-fmi-course/authorizer"
	"github.com/nikola-enter21/devops-fmi-course/gateway"
	"github.com/nikola-enter21/devops-fmi-course/logging"
	"github.com/nikola-enter21/devops-fmi-course/service/db/repo"
	"google.golang.org/grpc"
	health "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const readinessDrainDelay = 5 * time.Second

var (
	log = logging.MustNewLogger()
)

type Server struct {
	user.UnimplementedUserServiceServer

	Authorizer     authorizer.Authorizer
	UserRepository repo.UserRepository
}

func (s *Server) Serve(signalCtx context.Context, httpPort, grpcPort string) {
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
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(user.UserService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen", "error", err)
	}

	go gateway.Serve(signalCtx, ":"+httpPort, fmt.Sprintf("localhost:%s", grpcPort))

	go func() {
		log.Infow("gRPC server listening", "port", ":"+grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server exited", "error", err)
		}
	}()

	// wait for sigterm/sigint
	<-signalCtx.Done()

	// signal the LB that we are shutting down so it stops sending new traffic
	healthServer.SetServingStatus(user.UserService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_NOT_SERVING)

	// give some time for LB state to propagate
	time.Sleep(readinessDrainDelay)

	grpcServer.GracefulStop()
	log.Infow("gRPC server stopped cleanly")
}
