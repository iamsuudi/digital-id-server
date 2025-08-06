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
	q := repository.New(dbConn)

	log.Println("🌱 Seeding data...")
	
	// seedUsers(ctx, q)
	// seedLocations(ctx, q)
	seedResidents(ctx, q)
}
