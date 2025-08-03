package main

import (
	"context"

	"digital-id-server/internal/repository"
)

func seedResidents(ctx context.Context, queries *repository.Queries) {
	// Seed cities
	// city, err := queries.CreateCity(ctx, "Sample City")
	// if err != nil {
	// 	log.Fatalf("Failed to seed region: %v", err)
	// }

	// Seed kebeles
	// kebele, err := queries.CreateKebele(ctx, repository.CreateKebeleParams{
	// 	Name:   "Sample Region",
	// 	CityID: city.ID,
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to seed city: %v", err)
	// }
	
	// fmt.Println(kebele)

	// Seed addresses
	// addressID, err := queries.CreateAddress(ctx, repository.CreateAddressParams{
	// 	HouseNumber: "123",
	// 	KebeleID:    kebele.ID,
	// 	CityID:      city.ID,
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to seed address: %v", err)
	// }

	// Seed residents
	// residentID, err := queries.CreateResident(ctx, repository.CreateResidentParams{
	// 	Email:           "john.doe@example.com",
	// 	FirstName:       "John",
	// 	SecondName:      "Middle",
	// 	LastName:        "Doe",
	// 	BirthDate:       time.Now().AddDate(-30, 0, 0),
	// 	Gender:          "MALE",
	// 	Phone:           "+1234567890",
	// 	MaritalStatus:   "SINGLE",
	// 	Religion:        "NONE",
	// 	LanguagesSpoken: "English, Spanish",
	// 	AddressID:       &addressID,
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to seed resident: %v", err)
	// }

	// log.Printf("✅ Database seeded successfully. \nSample IDs: \nRegion=%s, \nCity=%s, \nAddress=%s, \nResident=%s",
	// 	kebele.ID, city.ID, addressID, residentID)
}
