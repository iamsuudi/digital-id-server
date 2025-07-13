package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/iamsuudi/digital-id-server/shared/config"
	db "github.com/iamsuudi/digital-id-server/database"
)

func main() {
	_ = godotenv.Load(".env")

	config.Load()
	conn := db.Connect()
	defer conn.Close()

	log.Println("🛠 Starting worker...")
	for {
		// Example placeholder task
		log.Println("💼 Worker is running background job...")
		time.Sleep(10 * time.Second)
	}
}
