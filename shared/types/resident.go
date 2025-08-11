package types

import (
	"time"
)

type ResidentPayload struct {
	Email      string    `form:"email"`
	FirstName  string    `form:"first_name" binding:"required"`
	SecondName string    `form:"second_name" binding:"required"`
	LastName   string    `form:"last_name" binding:"required"`
	BirthDate  time.Time `form:"birth_date"`
	Gender     string    `form:"gender"`
	Phone      string    `form:"phone"`

	MaritalStatus   string   `form:"marital_status"`
	Religion        string   `form:"religion"`
	Ethnicity       string   `form:"ethnicity"`
	Disability      string   `form:"disability"`
	EducationLevel  string   `form:"education_level"`
	NationalID      string   `form:"national_id"`
	LanguagesSpoken []string `form:"languages_spoken"`
	BloodType       string   `form:"blood_type"`

	EmploymentStatus string `form:"employment_status"`
	Occupation       string `form:"occupation"`
	EmploymentType   string `form:"employment_type"`
	EmployerName     string `form:"employer_name"`
	WorkAddress      string `form:"work_address"`

	EmergencyName     string `form:"emergency_contact_name"`
	EmergencyRelation string `form:"emergency_contact_relation"`
	EmergencyPhone    string `form:"emergency_contact_phone"`
	EmergencyEmail    string `form:"emergency_contact_email"`

	HouseNumber string `form:"house_number"`
	KebeleID    string `form:"kebele"`
	SubCityID   string `form:"subcity"`
	CityID      string `form:"city"`
}

type ResidentEmploymentPayload struct {
	ID               string `form:"id"`
	EmploymentStatus string `form:"status"`
	Occupation       string `form:"occupation"`
	EmploymentType   string `form:"type"`
	EmployerName     string `form:"employer_name"`
	WorkAddress      string `form:"work_address"`
}

type ResidentEmergencyPayload struct {
	ID                string `form:"id"`
	EmergencyName     string `form:"emergency_contact_name"`
	EmergencyRelation string `form:"emergency_contact_relation"`
	EmergencyPhone    string `form:"emergency_contact_phone"`
	EmergencyEmail    string `form:"emergency_contact_email"`
}

type ResidentPersonalPayload struct {
	ID         string    `form:"id"`
	Email      string    `form:"email"`
	FirstName  string    `form:"first_name" binding:"required"`
	SecondName string    `form:"second_name" binding:"required"`
	LastName   string    `form:"last_name" binding:"required"`
	BirthDate  time.Time `form:"birth_date"`
	Gender     string    `form:"gender"`
	Phone      string    `form:"phone"`
}

type ResidentAddressPayload struct {
	ID          string `form:"id"`
	HouseNumber string `form:"house_number"`
	KebeleID    string `form:"kebele_id"`
	SubCityID   string `form:"subcity_id"`
	CityID      string `form:"city_id"`
}

type ResidentAdditionalPayload struct {
	ID              string   `form:"id"`
	MaritalStatus   string   `form:"marital_status"`
	Religion        string   `form:"religion"`
	Ethnicity       string   `form:"ethnicity"`
	Disability      string   `form:"disability"`
	EducationLevel  string   `form:"education_level"`
	NationalID      string   `form:"national_id"`
	LanguagesSpoken []string `form:"languages_spoken"`
}
