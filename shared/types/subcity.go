package types

import "github.com/google/uuid"

type SubCityInput struct {
	Name      string     `json:"name" binding:"required"`
	CityID    uuid.UUID  `json:"city_id" binding:"required"`
	ManagerID *uuid.UUID `json:"manager_id"`
}
