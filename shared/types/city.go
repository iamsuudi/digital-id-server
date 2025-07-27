package types

import "github.com/google/uuid"

type CityInput struct {
	Name    string     `json:"name" binding:"required"`
	AdminID *uuid.UUID `json:"admin_id"`
}
