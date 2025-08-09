package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"digital-id-server/internal/repository"
	"digital-id-server/shared/types"

	"github.com/jackc/pgx/v5/pgxpool"
)

func seedResidents(ctx context.Context, db *pgxpool.Pool, q *repository.Queries) {
	// Delete existing residents
	if err := q.DeleteAllResidents(ctx); err != nil {
		log.Fatalf("Failed to delete residents: %v", err)
	}

	// Read JSON file
	defaultPath := "residents.json"
	if p := os.Getenv("RESIDENTS_FILE"); p != "" {
		defaultPath = p
	}
	filePath := flag.String("file", defaultPath, "path to residents.json")
	flag.Parse()
	data, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("Failed to read residents JSON: %v", err)
	}

	// Unmarshal JSON into slice
	var residents []types.ResidentPayload
	if err := json.Unmarshal(data, &residents); err != nil {
		log.Fatalf("Failed to unmarshal residents: %v", err)
	}

	// Insert each resident
	for _, r := range residents {
		tx, err := db.Begin(ctx)
		if err != nil {
			fmt.Printf("Couldn't create transaction: %v\n", err.Error())
			break
		}
		defer tx.Rollback(ctx)

		qtx := q.WithTx(tx)

		// 1. Create resident
		residentID, err := qtx.CreateResident(ctx, repository.CreateResidentParams{
			Email:      r.Email,
			FirstName:  r.FirstName,
			SecondName: r.SecondName,
			LastName:   r.LastName,
			BirthDate:  r.BirthDate,
			Gender:     r.Gender,
			Phone:      r.Phone,
		})
		if err != nil {
			fmt.Printf("Couldn't create %s: %v\n", r.FirstName, err.Error())
			break
		}

		// 2. Create biometric
		_, err = qtx.CreateBiometric(ctx, repository.CreateBiometricParams{
			ResidentID: residentID,
			BloodType:  r.BloodType,
			FaceUrl:    "face.png",
		})
		if err != nil {
			fmt.Printf("Couldn't biometric: %v\n", err.Error())
			break
		}

		// 3. Get or create address
		location, err := qtx.GetRandomLocation(ctx)
		if err != nil {
			fmt.Printf("Couldn't get random location: %v\n", err.Error())
			break
		}
		addr, err := qtx.GetAddressByLocations(ctx, repository.GetAddressByLocationsParams{
			HouseNumber: r.HouseNumber,
			KebeleID:    location.KebeleID,
			SubcityID:   location.SubcityID,
			CityID:      location.CityID,
		})
		if err != nil {
			newAddr, err := qtx.CreateAddress(ctx, repository.CreateAddressParams{
				HouseNumber: r.HouseNumber,
				KebeleID:    location.KebeleID,
				SubcityID:   location.SubcityID,
				CityID:      location.CityID,
			})
			if err != nil {
				fmt.Printf("Couldn't create new address: %v\n", err.Error())
				break
			}

			err = qtx.UpdateResidentAddress(ctx, repository.UpdateResidentAddressParams{
				AddressID: &(newAddr.ID),
				ID:        residentID,
			})
			if err != nil {
				fmt.Printf("Couldn't assign address to resident: %v\n", err.Error())
				break
			}
		} else {
			err = qtx.UpdateResidentAddress(ctx, repository.UpdateResidentAddressParams{
				AddressID: &addr.ID,
				ID:        residentID,
			})
			if err != nil {
				fmt.Printf("Couldn't assign address to resident: %v\n", err.Error())
				break
			}
		}

		// 4. Create document
		for _, doc := range []string{"doc1.png", "doc2.png"} {
			_, err = qtx.CreateDocument(ctx, repository.CreateDocumentParams{
				ResidentID: residentID,
				Url:        doc,
				Status:     "Pending",
			})
			if err != nil {
				fmt.Printf("Couldn't create ddocument: %v\n", err.Error())
				break
			}
		}

		// 5. Create employment
		_, err = qtx.CreateEmployment(ctx, repository.CreateEmploymentParams{
			ResidentID:   residentID,
			Status:       r.EmploymentStatus,
			Occupation:   &r.Occupation,
			EmployerName: &r.EmployerName,
			WorkAddress:  &r.WorkAddress,
			Type:         &r.EmploymentType,
		})
		if err != nil {
			fmt.Printf("Couldn't create employment: %v\n", err.Error())
			break
		}

		// 6. Create emergency
		_, err = qtx.CreateEmergencyContact(ctx, repository.CreateEmergencyContactParams{
			ResidentID: residentID,
			Name:       r.EmergencyName,
			Relation:   r.EmergencyRelation,
			Phone:      r.EmergencyPhone,
			Email:      &r.EmergencyEmail,
		})
		if err != nil {
			fmt.Printf("Couldn't create emergency: %v\n", err.Error())
			break
		}

		_, err = qtx.CreateAdditional(ctx, repository.CreateAdditionalParams{
			ResidentID:      residentID,
			Religion:        &r.Religion,
			Ethnicity:       &r.Ethnicity,
			NationalID:      &r.NationalID,
			Disability:      &r.Disability,
			EducationLevel:  &r.EducationLevel,
			MaritalStatus:   &r.MaritalStatus,
			LanguagesSpoken: r.LanguagesSpoken,
		})

		err = tx.Commit(ctx)
		if err != nil {
			fmt.Printf("Couldn't commit transaction: %v\n", err.Error())
			break
		}

		fmt.Printf("Resident %s registered\n", r.FirstName)
	}

	fmt.Println("✅ Residents seeded successfully")
}
