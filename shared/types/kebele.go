package types

import "github.com/google/uuid"

type KebeleInput struct {
	Name        string     `json:"name" binding:"required"`
	CityID      uuid.UUID  `json:"city_id" binding:"required"`
	SubCityID   uuid.UUID  `json:"subcity_id" binding:"required"`
	ExecutiveID *uuid.UUID `json:"executive_id"`
	Lat         *float64   `json:"lat"`
	Lon         *float64   `json:"lon"`
}
