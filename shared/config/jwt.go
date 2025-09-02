package config

import (
	"log"
	"os"
)

func GetJwtSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatalf("\n❌ Missing JWT_SECRET environment variable")
	}
	return secret
}
