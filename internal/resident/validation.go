package resident

import (
	"errors"
	"time"
)

func ValidateRegisterInput(input RegisterResidentInput) error {
	if input.Email == "" {
		return errors.New("email is required")
	}
	if input.FirstName == "" || input.LastName == "" {
		return errors.New("first and last names are required")
	}
	if input.BirthDate.IsZero() || input.BirthDate.After(time.Now()) {
		return errors.New("invalid birth date")
	}
	if input.Gender == "" {
		return errors.New("gender is required")
	}
	// ... add more as needed
	return nil
}
