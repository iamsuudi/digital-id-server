package main

import (
	"context"
	"log"

	"github.com/iamsuudi/digital-id-server/database"
	"github.com/iamsuudi/digital-id-server/internal/repository"
	"github.com/iamsuudi/digital-id-server/shared/config"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	config.Load()
	dbConn := database.Connect()
	defer dbConn.Close()

	ctx := context.Background()

	// Initialize sqlc Queries
	queries := repository.New(dbConn)

	log.Println("🌱 Seeding data...")

	// Seed users
	seedUsers(ctx, queries)
}
