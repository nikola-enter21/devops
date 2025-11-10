package config

import (
	"fmt"
	"os"
)

func Env(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func PostgresDSN() string {
	host := Env("DB_HOST", "postgres")
	port := Env("DB_PORT", "5432")
	usern := Env("DB_USER", "postgres")
	pass := Env("DB_PASSWORD", "postgres")
	dbname := Env("DB_NAME", "postgres")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, usern, pass, dbname,
	)

	return connStr
}
