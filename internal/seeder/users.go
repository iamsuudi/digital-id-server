package seeder

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func seedUsers(ctx context.Context, db *pgxpool.Pool) error {
	users := []struct {
		Name  string
		Email string
	}{
		{"Alice", "alice@example.com"},
		{"Bob", "bob@example.com"},
		{"Charlie", "charlie@example.com"},
	}

	for _, user := range users {
		_, err := db.Exec(ctx, `
			INSERT INTO users (name, email)
			VALUES ($1, $2)
			ON CONFLICT (email) DO NOTHING
		`, user.Name, user.Email)

		if err != nil {
			log.Printf("❌ Failed to seed user %s: %v\n", user.Email, err)
			return fmt.Errorf("failed to seed users: %w", err)
		}
	}

	log.Println("✅ Users seeded successfully.")
	return nil
}
