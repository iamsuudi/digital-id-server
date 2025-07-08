package main

import (
	"log"
	"time"

	"github.com/iamsuudi/digital-id-server/db"
	"github.com/iamsuudi/digital-id-server/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	cfg := config.Load()
	conn := db.Connect(cfg.DB)
	defer conn.Close()

	log.Println("🛠 Starting worker...")
	for {
		// Example placeholder task
		log.Println("💼 Worker is running background job...")
		time.Sleep(10 * time.Second)
	}
}
