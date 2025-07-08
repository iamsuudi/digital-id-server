package main

import (
	"context"
	"log"

	"github.com/iamsuudi/digital-id-server/db"
	"github.com/iamsuudi/digital-id-server/internal/config"
	"github.com/iamsuudi/digital-id-server/internal/seeder"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	cfg := config.Load()
	conn := db.Connect(cfg.DB)
	defer conn.Close()

	log.Println("🌱 Seeding data...")
	if err := seeder.SeedInitialUsers(context.Background(), conn); err != nil {
		log.Fatalf("❌ Failed to seed data: %v", err)
	}
	log.Println("✅ Seeding done.")
}
