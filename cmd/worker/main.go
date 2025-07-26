package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"digital-id-server/database"
	"digital-id-server/shared/config"
)

func main() {
	_ = godotenv.Load(".env")

	config.Load()
	conn := database.Connect()
	defer conn.Close()

	log.Println("🛠 Starting worker...")
	for {
		// Example placeholder task
		log.Println("💼 Worker is running background job...")
		time.Sleep(10 * time.Second)
	}
}
