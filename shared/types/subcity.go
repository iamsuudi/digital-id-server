package types

import "github.com/google/uuid"

type SubCityInput struct {
	Name string
	CityId uuid.UUID
}
