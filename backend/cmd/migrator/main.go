package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/nikola-enter21/devops-fmi-course/config"
	"github.com/pressly/goose/v3"
)

const migrationsDir = "migrations"

func main() {
	db, err := sql.Open("postgres", config.PostgresDSN())
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}()

	if err := goose.UpContext(context.Background(), db, migrationsDir); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("Migrations applied successfully")
}
