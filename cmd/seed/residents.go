package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"digital-id-server/internal/repository"
)

func seedResidents(ctx context.Context, q *repository.Queries) {
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
	var residents []repository.CreateResidentParams
	if err := json.Unmarshal(data, &residents); err != nil {
		log.Fatalf("Failed to unmarshal residents: %v", err)
	}

	// Insert each resident
	for _, r := range residents {
		if _, err := q.CreateResident(ctx, r); err != nil {
			log.Fatalf("Failed to insert resident: %v", err)
		}
		fmt.Println("Inserted resident:", r.FirstName)
	}

	fmt.Println("✅ Residents seeded successfully")
}
