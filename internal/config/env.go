package config

import (
	"fmt"
	"log"
	"os"

	"github.com/iamsuudi/digital-id-server/db"
	"github.com/joho/godotenv"
)

type Config struct {
	DB db.DBConfig
}

func Load() Config {
	fmt.Println("🔧 Loading configuration...")
	// Load .env.mk first (Makefile-style values)
	if err := godotenv.Overload(".env.mk"); err != nil {
		log.Println("⚠️ .env.mk not found or unreadable (optional for Go)")
	}
	
	// Load .env next — allows overrides and complex keys
	if err := godotenv.Overload(".env"); err != nil {
		log.Println("⚠️ .env not found or unreadable (optional)")
	}
	fmt.Println("🔧 Environment variables loaded successfully")

	return Config{
		DB: db.DBConfig{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
		},
	}
}
