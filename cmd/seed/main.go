package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"digital-id-server/database"
	"digital-id-server/shared/config"
	"digital-id-server/internal/repository"
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

	seedUsers(ctx, queries)
	seedLocations(ctx, queries)
}
