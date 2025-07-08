package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

// validateEnv checks that required env vars are present
func validateEnv(cfg DBConfig) {
	missing := []string{}
	if cfg.User == "" {
		missing = append(missing, "DB_USER")
	}
	/* if cfg.Password == "" {
		missing = append(missing, "DB_PASSWORD")
	} */
	if cfg.Host == "" {
		missing = append(missing, "DB_HOST")
	}
	if cfg.Port == "" {
		missing = append(missing, "DB_PORT")
	}
	if cfg.Name == "" {
		missing = append(missing, "DB_NAME")
	}
	if len(missing) > 0 {
		log.Fatalf("\n❌ Missing required DB environment variables: %v", missing)
	}
}

func Connect(cfg DBConfig) *pgxpool.Pool {
	validateEnv(cfg)
	
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
