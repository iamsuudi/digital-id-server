package seeder

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SeedInitialUsers seeds basic users (you can call more seeders here too)
func SeedInitialUsers(ctx context.Context, db *pgxpool.Pool) error {
	log.Println("👤 Seeding users...")

	if err := seedUsers(ctx, db); err != nil {
		return err
	}

	// You can add more:
	// if err := seedPosts(ctx, db); err != nil {
	//     return err
	// }

	return nil
}
