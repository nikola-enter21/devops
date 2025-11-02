package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/nikola-enter21/devops-fmi-course/api/gen/go/user/v1"
	"github.com/nikola-enter21/devops-fmi-course/logging"
)

var (
	log = logging.MustNewLogger()
)

type UserServiceServer struct {
	user.UnimplementedUserServiceServer
}

func (s *UserServiceServer) Healthz(ctx context.Context, _ *user.HealthzRequest) (*user.HealthzResponse, error) {
	log.Infow("health check")
	return &user.HealthzResponse{
		Status: "ok",
	}, nil
}

func (s *UserServiceServer) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	log.Infow("login attempt", "username", req.Username)
	return &user.LoginResponse{
		Token: "Login successful",
	}, nil
}

func (s *UserServiceServer) Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	log.Infow("user registration", "username", req.Username, "email", req.Email)
	return &user.RegisterResponse{
		Message: "User registered",
	}, nil
}

func (s *UserServiceServer) CheckDatabase(ctx context.Context, _ *user.CheckDatabaseRequest) (*user.CheckDatabaseResponse, error) {
	host := getEnv("DB_HOST", "postgres")
	port := getEnv("DB_PORT", "5432")
	usern := getEnv("DB_USER", "postgres")
	pass := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "postgres")

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

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
