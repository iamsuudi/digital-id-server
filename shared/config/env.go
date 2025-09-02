package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"digital-id-server/shared/types"
)

func Load() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	// Files to load, in order
	files := []string{
		".env",
		".env." + env,
		".env." + env + ".mk",
	}

	for _, file := range files {
		if err := godotenv.Overload(file); err != nil {
			log.Printf("⚠️ %s not found or unreadable.\n", file)
		} else {
			fmt.Printf("✅ Loaded %s\n", file)
		}
	}

	fmt.Printf("🔧 Environment variables loaded for %s\n", env)
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
