package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"digital-id-server/internal/repository"
)

type Kebele struct {
	Name string
	Lat  float64
	Lon  float64
}

type SubCity struct {
	Name    string
	Lat     float64
	Lon     float64
	Kebeles []Kebele
}

type City struct {
	Name      string
	Lat       float64
	Lon       float64
	SubCities []SubCity
}

var Data = []City{
	{
		Name: "Adama",
		Lat:  8.5263, Lon: 39.2583,
		SubCities: []SubCity{
			{
				Name: "Subcity-01", Lat: 8.520, Lon: 39.270,
				Kebeles: []Kebele{
					{"Kebele 01", 8.521, 39.271},
					{"Kebele 02", 8.522, 39.268},
					{"Kebele 03", 8.519, 39.269},
				},
			},
			{
				Name: "Subcity-02", Lat: 8.510, Lon: 39.265,
				Kebeles: []Kebele{
					{"Kebele 01", 8.511, 39.266},
					{"Kebele 02", 8.509, 39.264},
					{"Kebele 03", 8.512, 39.265},
				},
			},
			// ... similar subcities 3–5
		},
	},
	{
		Name: "Jimma",
		Lat:  7.6753, Lon: 36.8373, // approximate
		SubCities: []SubCity{
			{
				Name: "Subcity-01", Lat: 7.670, Lon: 36.839,
				Kebeles: []Kebele{
					{"Agaro Sefer", 7.671, 36.840},
					{"Baker Sefer", 7.669, 36.838},
					{"Merkato", 7.672, 36.837},
				},
			},
			{
				Name: "Subcity-02", Lat: 7.662, Lon: 36.830,
				Kebeles: []Kebele{
					{"Higher 01", 7.663, 36.831},
					{"Higher 02", 7.661, 36.829},
					{"Higher 03", 7.664, 36.830},
				},
			},
			// Subcity-03, 04 same pattern
		},
	},
	{
		Name: "Shashamane",
		Lat:  7.1964, Lon: 38.5977, // approximate
		SubCities: []SubCity{
			{
				Name: "Subcity-02", Lat: 7.193, Lon: 38.590,
				Kebeles: []Kebele{
					{"Kebele 01", 7.194, 38.591},
					{"Kebele 02", 7.192, 38.589},
				},
			},
			{
				Name: "Subcity-03", Lat: 7.198, Lon: 38.593,
				Kebeles: []Kebele{
					{"Kebele 01", 7.199, 38.594},
					{"Kebele 02", 7.197, 38.592},
				},
			},
			{
				Name: "Subcity-04", Lat: 7.195, Lon: 38.599,
				Kebeles: []Kebele{
					{"Kebele 01", 7.196, 38.600},
					{"Kebele 02", 7.194, 38.598},
				},
			},
		},
	},
	{
		Name: "Bishoftu",
		Lat:  8.7481, Lon: 38.97868,
		SubCities: []SubCity{
			{
				Name: "Subcity-02", Lat: 8.745, Lon: 38.976,
				Kebeles: []Kebele{
					{"Kebele 01", 8.746, 38.977},
					{"Kebele 02", 8.744, 38.975},
				},
			},
			{
				Name: "Subcity-03", Lat: 8.750, Lon: 38.982,
				Kebeles: []Kebele{
					{"Kebele 01", 8.751, 38.983},
					{"Kebele 02", 8.749, 38.981},
				},
			},
		},
	},
	{
		Name: "Nekemte",
		Lat:  9.0881, Lon: 36.5472, // approximate
		SubCities: []SubCity{
			{
				Name: "Subcity-02", Lat: 9.090, Lon: 36.543,
				Kebeles: []Kebele{
					{"Kebele 01", 9.091, 36.544},
					{"Kebele 02", 9.089, 36.542},
				},
			},
			{
				Name: "Subcity-03", Lat: 9.087, Lon: 36.546,
				Kebeles: []Kebele{
					{"Kebele 01", 9.088, 36.547},
					{"Kebele 02", 9.086, 36.545},
				},
			},
		},
	},
	{
		Name: "Ambo",
		Lat:  8.9821, Lon: 37.8582, // approximate
		SubCities: []SubCity{
			{
				Name: "Subcity-02", Lat: 8.985, Lon: 37.860,
				Kebeles: []Kebele{
					{"Kebele 01", 8.986, 37.861},
					{"Kebele 02", 8.984, 37.859},
				},
			},
		},
	},
	{
		Name: "Asella",
		Lat:  7.9565, Lon: 39.1321, // approximate
		SubCities: []SubCity{
			{
				Name: "Subcity-01", Lat: 7.987, Lon: 39.125,
				Kebeles: []Kebele{
					{"Kebele 01", 7.988, 39.126},
					{"Kebele 02", 7.986, 39.124},
				},
			},
		},
	},
	{
		Name: "Burayu",
		Lat:  9.0341, Lon: 38.6619, // approximate
		SubCities: []SubCity{
			{
				Name: "Subcity-01", Lat: 9.051, Lon: 38.652,
				Kebeles: []Kebele{
					{"Kebele 01", 9.052, 38.653},
					{"Kebele 02", 9.050, 38.651},
					{"Kebele 03", 9.053, 38.650},
				},
			},
		},
	},
	{
		Name: "Sebeta",
		Lat:  8.9112, Lon: 38.6268, // approximate
		SubCities: []SubCity{
			{
				Name: "Subcity-01", Lat: 8.852, Lon: 38.652,
				Kebeles: []Kebele{
					{"Kebele 01", 8.853, 38.653},
					{"Kebele 02", 8.851, 38.651},
				},
			},
			{
				Name: "Subcity-02", Lat: 8.848, Lon: 38.648,
				Kebeles: []Kebele{
					{"Kebele 01", 8.849, 38.649},
					{"Kebele 02", 8.847, 38.647},
					{"Kebele 03", 8.850, 38.646},
				},
			},
		},
	},
	{
		Name: "Gimbi",
		Lat:  9.1773, Lon: 35.8339, // approximate
		SubCities: []SubCity{
			{
				Name: "Subcity-01", Lat: 9.165, Lon: 36.505,
				Kebeles: []Kebele{
					{"Kebele 01", 9.166, 36.506},
					{"Kebele 02", 9.164, 36.504},
				},
			},
			{
				Name: "Subcity-02", Lat: 9.160, Lon: 36.500,
				Kebeles: []Kebele{
					{"Kebele 01", 9.161, 36.501},
					{"Kebele 02", 9.159, 36.499},
				},
			},
		},
	},
}

func seedLocations(ctx context.Context, queries *repository.Queries) {
	start := time.Now()

	for _, c := range Data {
		city, err := queries.CreateCity(ctx, repository.CreateCityParams{
			Name: c.Name,
			Lat:  &c.Lat,
			Lon:  &c.Lon,
		})
		if err != nil {
			log.Fatalf("Failed to create city: %v", c.Name)
		}
		log.Printf("✅ City created: %s.", city.Name)

		for _, sc := range c.SubCities {
			subCity, err := queries.CreateSubCity(ctx, repository.CreateSubCityParams{
				Name:   sc.Name,
				Lat:    &c.Lat,
				Lon:    &c.Lon,
				CityID: city.ID,
			})
			if err != nil {
				log.Fatalf("Failed to create sub-city: %v", sc.Name)
			}
			log.Printf("✅ Sub-city created: %s.", subCity.Name)

			for _, k := range sc.Kebeles {
				kebele, err := queries.CreateKebele(ctx, repository.CreateKebeleParams{
					Name:      k.Name,
					Lat:       &c.Lat,
					Lon:       &c.Lon,
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
