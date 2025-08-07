package types

import (
	"time"

	"github.com/google/uuid"
)

type ResidentPayload struct {
	Email             string    `form:"email"`
	FirstName         string    `form:"first_name" binding:"required"`
	SecondName        string    `form:"second_name" binding:"required"`
	LastName          string    `form:"last_name" binding:"required"`
	BirthDate         time.Time `form:"birth_date"`
	Gender            string    `form:"gender"`
	Phone             string    `form:"phone"`
	MaritalStatus     string    `form:"marital_status"`
	Religion          string    `form:"religion"`
	Ethnicity         string    `form:"ethnicity"`
	Disability        string    `form:"disability"`
	EducationLevel    string    `form:"educational_level"`
	LanguagesSpoken   []string  `form:"languages_spoken"`
	BloodType         string    `form:"blood_type"`
	Face              string    `form:"face"`
	DocumentType      string    `form:"document_type"`
	DocumentNumber    string    `form:"document_number"`
	EmploymentStatus  string    `form:"document_status"`
	Occupation        string    `form:"occupation"`
	EmployerName      string    `form:"employer_name"`
	WorkAddress       string    `form:"work_address"`
	EmergencyName     string    `form:"emergency_contact_name"`
	EmergencyRelation string    `form:"emergency_contact_relation"`
	EmergencyPhone    string    `form:"emergency_contact_phone"`
	HouseNumber       string    `form:"house_number"`
	KebeleID          uuid.UUID `form:"kebele_id"`
	SubCityID         uuid.UUID `form:"subcity_id"`
	CityID            uuid.UUID `form:"city_id"`
}
