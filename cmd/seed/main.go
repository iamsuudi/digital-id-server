package main

import (
	"context"
	"log"
	"time"

	"github.com/iamsuudi/digital-id-server/database"
	"github.com/iamsuudi/digital-id-server/internal/repository"
	"github.com/iamsuudi/digital-id-server/shared/config"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load(".env")

	config.Load()
	dbConn := database.Connect()
	defer dbConn.Close()

	ctx := context.Background()

	// Initialize sqlc Queries
	queries := repository.New(dbConn)

	log.Println("🌱 Seeding data...")

	// Seed cities
	city, err := queries.CreateCity(ctx, repository.CreateCityParams{
		Name:      "Sample City",
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Fatalf("Failed to seed region: %v", err)
	}

	// Seed regions
	kebele, err := queries.CreateKebele(ctx, repository.CreateKebeleParams{
		Name:      "Sample Region",
		CityID:    city.ID,
		CreatedAt: time.Now(),
	})

	if err != nil {
		log.Fatalf("Failed to seed city: %v", err)
	}

	// Seed addresses
	addressID, err := queries.CreateAddress(ctx, repository.CreateAddressParams{
		HouseNumber: "123",
		KebeleID:    kebele.ID,
		CityID:      city.ID,
	})
	if err != nil {
		log.Fatalf("Failed to seed address: %v", err)
	}

	// Seed residents
	residentID, err := queries.CreateResident(ctx, repository.CreateResidentParams{
		Email:           "john.doe@example.com",
		FirstName:       "John",
		SecondName:      "Middle",
		LastName:        "Doe",
		BirthDate:       time.Now().AddDate(-30, 0, 0),
		Gender:          "MALE",
		Phone:           "+1234567890",
		MaritalStatus:   "SINGLE",
		Religion:        "NONE",
		LanguagesSpoken: "English, Spanish",
		AddressID:       &addressID,
	})
	if err != nil {
		log.Fatalf("Failed to seed resident: %v", err)
	}

	// Seed users (e.g., Super Admin)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	userID, err := queries.CreateUser(ctx, repository.CreateUserParams{
		FirstName:    "Admin",
		SecondName:   "Super",
		LastName:     "User",
		Email:        "admin@example.com",
		Phone:        "+1234567890",
		PasswordHash: string(hashedPassword),
		Role:         "SUPERADMIN",
	})
	if err != nil {
		log.Fatalf("Failed to seed user: %v", err)
	}

	log.Printf("✅ Database seeded successfully. \nSample IDs: \nRegion=%s, \nCity=%s, \nAddress=%s, \nResident=%s, \nUser=%s",
		kebele.ID, city.ID, addressID, residentID, userID)
}
