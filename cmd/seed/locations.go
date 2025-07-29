package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"digital-id-server/internal/repository"
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
			{Name: "Adama Subcity-01", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}, {Name: "Kebele 03"}}},
			{Name: "Adama Subcity-02", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}, {Name: "Kebele 03"}}},
			{Name: "Adama Subcity-03", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}}},
			{Name: "Adama Subcity-04", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}, {Name: "Kebele 03"}}},
			{Name: "Adama Subcity-05", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}}},
		},
	},
	{
		Name: "Jimma",
		SubCities: []SubCity{
			{Name: "Jimma Subcity-01", Kebeles: []Kebele{{Name: "Agaro Sefer"}, {Name: "Baker Sefer"}, {Name: "Merkato"}}},
			{Name: "Jimma Subcity-02", Kebeles: []Kebele{{Name: "Higher 01"}, {Name: "Higher 02"}, {Name: "Higher 03"}}},
			{Name: "Jimma Subcity-03", Kebeles: []Kebele{{Name: "Town 01"}, {Name: "Town 02"}}},
			{Name: "Jimma Subcity-04", Kebeles: []Kebele{{Name: "Seka"}, {Name: "Mizan-Aman"}}},
		},
	},
	{
		Name: "Shashamane",
		SubCities: []SubCity{
			{Name: "Shashemene Subcity-01", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}}},
			{Name: "Shashemene Subcity-02", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}, {Name: "Kebele 03"}}},
			{Name: "Shashemene Subcity-03", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}}},
			{Name: "Shashemene Subcity-04", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}}},
		},
	},
	{
		Name: "Bishoftu",
		SubCities: []SubCity{
			{Name: "Bishoftu Subcity-01", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}, {Name: "Kebele 03"}}},
			{Name: "Bishoftu Subcity-02", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}}},
			{Name: "Bishoftu Subcity-03", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}}},
		},
	},
	{
		Name: "Nekemte",
		SubCities: []SubCity{
			{Name: "Nekemte Subcity-01", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}, {Name: "Kebele 03"}}},
			{Name: "Nekemte Subcity-02", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}}},
			{Name: "Nekemte Subcity-03", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}}},
		},
	},
	{
		Name: "Ambo",
		SubCities: []SubCity{
			{Name: "Ambo Subcity-01", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}}},
			{Name: "Ambo Subcity-02", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}, {Name: "Kebele 03"}}},
			{Name: "Ambo Subcity-03", Kebeles: []Kebele{{Name: "Kebele 01"}}},
		},
	},
	{
		Name: "Asella",
		SubCities: []SubCity{
			{Name: "Asella Subcity-01", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}}},
			{Name: "Asella Subcity-02", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}}},
			{Name: "Asella Subcity-03", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}}},
		},
	},
	{
		Name: "Burayu",
		SubCities: []SubCity{
			{Name: "Burayu Subcity-01", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}, {Name: "Kebele 03"}}},
			{Name: "Burayu Subcity-02", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}}},
		},
	},
	{
		Name: "Sebeta",
		SubCities: []SubCity{
			{Name: "Sebeta Subcity-01", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}}},
			{Name: "Sebeta Subcity-02", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}, {Name: "Kebele 03"}}},
			{Name: "Sebeta Subcity-03", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}}},
		},
	},
	{
		Name: "Gimbi",
		SubCities: []SubCity{
			{Name: "Gimbi Subcity-01", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}}},
			{Name: "Gimbi Subcity-02", Kebeles: []Kebele{{Name: "Kebele 01"}, {Name: "Kebele 02"}}},
		},
	},
}

func seedLocations(ctx context.Context, queries *repository.Queries) {
	start := time.Now()

	for _, c := range Data {
		city, err := queries.CreateCity(ctx, c.Name)
		if err != nil {
			log.Fatalf("Failed to create city: %v", c.Name)
		}
		log.Printf("✅ City created: %s.", city.Name)

		for _, sc := range c.SubCities {
			subCity, err := queries.CreateSubCity(ctx, repository.CreateSubCityParams{
				Name:   sc.Name,
				CityID: city.ID,
			})
			if err != nil {
				log.Fatalf("Failed to create sub-city: %v", sc.Name)
			}
			log.Printf("✅ Sub-city created: %s.", subCity.Name)

			for _, k := range sc.Kebeles {
				kebele, err := queries.CreateKebele(ctx, repository.CreateKebeleParams{
					Name:      k.Name,
					SubcityID: &(subCity.ID),
					CityID:    city.ID,
				})
				if err != nil {
					log.Fatalf("Failed to create kebele: %v", k.Name)
				}
				log.Printf("✅ Kebele created: %s.", kebele.Name)
			}
		}
	}

	elapsed := time.Since(start)

	fmt.Printf("\n✅ Locations seeded successfully. Took %s\n", elapsed)
}
