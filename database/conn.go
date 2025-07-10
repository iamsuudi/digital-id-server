package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/iamsuudi/digital-id-server/shared/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() *pgxpool.Pool {
	cfg := config.GetDatabaseConfig()
	
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("\n❌ Failed to create connection pool: %v", err)
	}

	// Ping the DB to ensure it's up
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("\n❌ Database not reachable: %v", err)
	}

	log.Println("✅ Connected to PostgreSQL successfully")
	return pool
}
