package types

import (
	"time"

	"github.com/google/uuid"
)

type ResidentInput struct {
	Email             string
	FirstName         string
	SecondName        string
	LastName          string
	BirthDate         time.Time
	Gender            string
	Phone             string
	MaritalStatus     string
	Religion          string
	Ethnicity         string
	Disability        string
	EducationLevel    string
	LanguagesSpoken   []string
	BloodType         string
	Face              string // this might later be a file or URL
	DocumentType      string
	DocumentNumber    string
	DocumentURL       string
	EmploymentStatus  string
	Occupation        string
	EmployerName      string
	WorkAddress       string
	EmergencyName     string
	EmergencyRelation string
	EmergencyPhone    string
	HouseNumber       string
	KebeleID          uuid.UUID
	CityID            uuid.UUID
}
