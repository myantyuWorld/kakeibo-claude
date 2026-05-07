package main

import (
	"log"

	"backend/migrations"
	"backend/pkg/config"
	"backend/pkg/postgres"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := postgres.NewDB(cfg.Postgres.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	if err := db.Migrate(migrations.FS); err != nil {
		db.Close()
		log.Fatalf("migration failed: %v", err)
	}
	db.Close()
	log.Println("migration completed")
}
