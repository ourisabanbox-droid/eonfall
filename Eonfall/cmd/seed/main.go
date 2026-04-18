package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"project-eonfall/internal/config"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config.Load: %v", err)
	}

	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("pgxpool.New: %v", err)
	}
	defer db.Close()

	if err := seedWorld(ctx, db); err != nil {
		log.Fatalf("seedWorld: %v", err)
	}

	log.Println("seed completed")
}
