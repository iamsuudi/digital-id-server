package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"digital-id-server/shared/types"
)

func Load() {
	fmt.Println("🔧 Loading configuration...")
	// Load .env.mk first (Makefile-style values)
	if err := godotenv.Overload(".env.mk"); err != nil {
		log.Println("⚠️ .env.mk not found or unreadable.")
	}

	// Load .env next — allows overrides and complex keys
	if err := godotenv.Overload(".env"); err != nil {
		log.Println("⚠️ .env not found or unreadable.")
	}
	fmt.Println("🔧 Environment variables loaded successfully")
}

func GetDatabaseConfig() types.DBConfig {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	missing := []string{}
	if user == "" {
		missing = append(missing, "DB_USER")
	}
	if password == "" {
		// missing = append(missing, "DB_PASSWORD")
	}
	if host == "" {
		missing = append(missing, "DB_HOST")
	}
	if port == "" {
		missing = append(missing, "DB_PORT")
	}
	if name == "" {
		missing = append(missing, "DB_NAME")
	}
	if len(missing) > 0 {
		log.Fatalf("\n❌ Missing required DB environment variables: %v", missing)
	}

	return types.DBConfig{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
		Name:     name,
	}
}

func GetJwtSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatalf("\n❌ Missing JWT_SECRET environment variable")
	}
	return secret
}
