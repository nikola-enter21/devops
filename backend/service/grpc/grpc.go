package grpc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/nikola-enter21/devops-fmi-course/api/gen/go/user/v1"
	"github.com/nikola-enter21/devops-fmi-course/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Healthz(ctx context.Context, _ *user.HealthzRequest) (*user.HealthzResponse, error) {
	log.Infow("health check")
	return &user.HealthzResponse{
		Status: "ok",
	}, nil
}

func (s *Server) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	log.Infow("login attempt", "username", req.Username)

	_, err := s.UserRepository.GetByID(ctx, 123)
	if err != nil {
		log.Errorw("user by id failed", "error", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &user.LoginResponse{
		Token: "Login successful",
	}, nil
}

func (s *Server) Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	log.Infow("user registration", "username", req.Username, "email", req.Email)
	return &user.RegisterResponse{
		Message: "User registered",
	}, nil
}

func (s *Server) CheckDatabase(ctx context.Context, _ *user.CheckDatabaseRequest) (*user.CheckDatabaseResponse, error) {
	host := config.Env("DB_HOST", "postgres")
	port := config.Env("DB_PORT", "5432")
	usern := config.Env("DB_USER", "postgres")
	pass := config.Env("DB_PASSWORD", "postgres")
	dbname := config.Env("DB_NAME", "postgres")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, usern, pass, dbname,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Errorw("failed to open DB connection", "error", err)
		return &user.CheckDatabaseResponse{
			DbStatus: err.Error(),
		}, nil
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			log.Errorw("failed to close DB connection", "error", cerr)
		}
	}()

	if err := db.PingContext(ctx); err != nil {
		log.Errorw("database unreachable", "error", err)
		return &user.CheckDatabaseResponse{
			DbStatus: "unreachable: " + err.Error(),
		}, nil
	}

	log.Infow("database connection successful",
		"host", host,
		"port", port,
		"user", usern,
		"dbname", dbname,
	)

	return &user.CheckDatabaseResponse{
		DbStatus: "ok",
	}, nil
}
