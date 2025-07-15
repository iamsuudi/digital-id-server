package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/iamsuudi/digital-id-server/internal/repository"
)

// Kebele is the smallest administrative unit.
type Kebele struct {
	Name string
}

// SubCity groups kebeles.
type SubCity struct {
	Name    string
	Kebeles []Kebele
}

// City groups sub-cities.
type City struct {
	Name      string
	SubCities []SubCity
}

// Data holds the 10 representative Oromia cities.
var Data = []City{
	{
		Name: "Adama",
		SubCities: []SubCity{
			{Name: "Adama Ketema-01", Kebeles: []Kebele{{Name: "01"}, {Name: "02"}, {Name: "03"}}},
			{Name: "Adama Ketema-02", Kebeles: []Kebele{{Name: "04"}, {Name: "05"}, {Name: "06"}}},
			{Name: "Adama Ketema-03", Kebeles: []Kebele{{Name: "07"}, {Name: "08"}}},
			{Name: "Adama Ketema-04", Kebeles: []Kebele{{Name: "09"}, {Name: "10"}, {Name: "11"}}},
			{Name: "Adama Ketema-05", Kebeles: []Kebele{{Name: "12"}, {Name: "13"}}},
		},
	},
	{
		Name: "Jimma",
		SubCities: []SubCity{
			{Name: "Jimma Ketema-01", Kebeles: []Kebele{{Name: "Agaro Sefer"}, {Name: "Baker Sefer"}, {Name: "Merkato"}}},
			{Name: "Jimma Ketema-02", Kebeles: []Kebele{{Name: "Higher 01"}, {Name: "Higher 02"}, {Name: "Higher 03"}}},
			{Name: "Jimma Ketema-03", Kebeles: []Kebele{{Name: "Town 01"}, {Name: "Town 02"}}},
			{Name: "Jimma Ketema-04", Kebeles: []Kebele{{Name: "Seka"}, {Name: "Mizan-Aman"}}},
		},
	},
	{
		Name: "Shashamane",
		SubCities: []SubCity{
			{Name: "Shashemene Ketema-01", Kebeles: []Kebele{{Name: "01"}, {Name: "02"}}},
			{Name: "Shashemene Ketema-02", Kebeles: []Kebele{{Name: "03"}, {Name: "04"}, {Name: "05"}}},
			{Name: "Shashemene Ketema-03", Kebeles: []Kebele{{Name: "06"}, {Name: "07"}}},
			{Name: "Shashemene Ketema-04", Kebeles: []Kebele{{Name: "08"}, {Name: "09"}}},
		},
	},
	{
		Name: "Bishoftu",
		SubCities: []SubCity{
			{Name: "Bishoftu Ketema-01", Kebeles: []Kebele{{Name: "01"}, {Name: "02"}, {Name: "03"}}},
			{Name: "Bishoftu Ketema-02", Kebeles: []Kebele{{Name: "04"}, {Name: "05"}}},
			{Name: "Bishoftu Ketema-03", Kebeles: []Kebele{{Name: "06"}, {Name: "07"}}},
		},
	},
	{
		Name: "Nekemte",
		SubCities: []SubCity{
			{Name: "Nekemte Ketema-01", Kebeles: []Kebele{{Name: "01"}, {Name: "02"}, {Name: "03"}}},
			{Name: "Nekemte Ketema-02", Kebeles: []Kebele{{Name: "04"}, {Name: "05"}}},
			{Name: "Nekemte Ketema-03", Kebeles: []Kebele{{Name: "06"}, {Name: "07"}}},
		},
	},
	{
		Name: "Ambo",
		SubCities: []SubCity{
			{Name: "Ambo Ketema-01", Kebeles: []Kebele{{Name: "01"}, {Name: "02"}}},
			{Name: "Ambo Ketema-02", Kebeles: []Kebele{{Name: "03"}, {Name: "04"}, {Name: "05"}}},
			{Name: "Ambo Ketema-03", Kebeles: []Kebele{{Name: "06"}}},
		},
	},
	{
		Name: "Asella",
		SubCities: []SubCity{
			{Name: "Asella Ketema-01", Kebeles: []Kebele{{Name: "01"}, {Name: "02"}}},
			{Name: "Asella Ketema-02", Kebeles: []Kebele{{Name: "03"}, {Name: "04"}}},
			{Name: "Asella Ketema-03", Kebeles: []Kebele{{Name: "05"}, {Name: "06"}}},
		},
	},
	{
		Name: "Burayu",
		SubCities: []SubCity{
			{Name: "Burayu Ketema-01", Kebeles: []Kebele{{Name: "01"}, {Name: "02"}, {Name: "03"}}},
			{Name: "Burayu Ketema-02", Kebeles: []Kebele{{Name: "04"}, {Name: "05"}}},
		},
	},
	{
		Name: "Sebeta",
		SubCities: []SubCity{
			{Name: "Sebeta Ketema-01", Kebeles: []Kebele{{Name: "01"}, {Name: "02"}}},
			{Name: "Sebeta Ketema-02", Kebeles: []Kebele{{Name: "03"}, {Name: "04"}, {Name: "05"}}},
			{Name: "Sebeta Ketema-03", Kebeles: []Kebele{{Name: "06"}, {Name: "07"}}},
		},
	},
	{
		Name: "Gimbi",
		SubCities: []SubCity{
			{Name: "Gimbi Ketema-01", Kebeles: []Kebele{{Name: "01"}, {Name: "02"}}},
			{Name: "Gimbi Ketema-02", Kebeles: []Kebele{{Name: "03"}, {Name: "04"}}},
		},
	},
}

func seedLocations(ctx context.Context, queries *repository.Queries) {
	start := time.Now()
	
	for _, c := range Data {
		_, err := queries.CreateCity(ctx, c.Name)
		if err != nil {
			log.Fatalf("Failed to create city: %v", c.Name)
		}
		// log.Printf("✅ City created: %s.", city.Name)
	}
	
	elapsed := time.Since(start)

	fmt.Printf("\n✅ %d cities seeded successfully. Took %s\n", len(Data), elapsed)
}