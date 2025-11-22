package grpc

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/nikola-enter21/devops-fmi-course/api/gen/go/user/v1"
	"github.com/nikola-enter21/devops-fmi-course/config"
)

func (s *Server) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	return &user.LoginResponse{
		Token: "Login successful",
	}, nil
}

func (s *Server) Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
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

	return &user.CheckDatabaseResponse{
		DbStatus: "ok",
	}, nil
}
