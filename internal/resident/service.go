package resident

import (
	"context"
	"fmt"

	"digital-id-server/internal/repository"
	"digital-id-server/shared/types"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	db *pgxpool.Pool
	q  *repository.Queries
}

func NewService(dbConn *pgxpool.Pool, dbQueries *repository.Queries) *Service {
	return &Service{db: dbConn, q: dbQueries}
}

func (s *Service) RegisterResident(ctx context.Context, input types.ResidentPayload, docsUrl []string, faceUrl string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	// 1. Create resident
	residentID, err := qtx.CreateResident(ctx, repository.CreateResidentParams{
		Email:      input.Email,
		FirstName:  input.FirstName,
		SecondName: input.SecondName,
		LastName:   input.LastName,
		BirthDate:  input.BirthDate,
		Gender:     input.Gender,
		Phone:      input.Phone,
	})
	if err != nil {
		return err
	}

	// 2. Create biometric
	_, err = qtx.CreateBiometric(ctx, repository.CreateBiometricParams{
		ResidentID: residentID,
		BloodType:  input.BloodType,
		FaceUrl:    faceUrl,
	})
	if err != nil {
		return err
	}
	fmt.Println("Biometric created")

	raw := input.KebeleID
	kebeleID, err := uuid.Parse(raw)
	if err != nil {
		return err
	}

	raw = input.SubCityID
	subcityID, err := uuid.Parse(raw)
	if err != nil {
		return err
	}

	raw = input.CityID
	cityID, err := uuid.Parse(raw)
	if err != nil {
		return err
	}

	// 3. Get or create address
	addr, err := qtx.GetAddressByLocations(ctx, repository.GetAddressByLocationsParams{
		HouseNumber: input.HouseNumber,
		KebeleID:    kebeleID,
		SubcityID:   subcityID,
		CityID:      cityID,
	})
	if err != nil {
		fmt.Println("Address not found")
		newAddr, err := qtx.CreateAddress(ctx, repository.CreateAddressParams{
			HouseNumber: input.HouseNumber,
			KebeleID:    kebeleID,
			SubcityID:   subcityID,
			CityID:      cityID,
		})
		if err != nil {
			return err
		}
		fmt.Println("New address created")
		err = qtx.UpdateResidentAddress(ctx, repository.UpdateResidentAddressParams{
			AddressID: &(newAddr.ID),
			ID:        residentID,
		})
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Address found")
		err = qtx.UpdateResidentAddress(ctx, repository.UpdateResidentAddressParams{
			AddressID: &addr.ID,
			ID:        residentID,
		})
		if err != nil {
			return err
		}
	}

	// 4. Create document
	for _, doc := range docsUrl {
		_, err = qtx.CreateDocument(ctx, repository.CreateDocumentParams{
			ResidentID: residentID,
			Url:        doc,
			Status:     "Pending",
		})
		if err != nil {
			return err
		}
	}
	fmt.Println("Documents created")

	// 5. Create employment
	_, err = qtx.CreateEmployment(ctx, repository.CreateEmploymentParams{
		ResidentID:   residentID,
		Status:       input.EmploymentStatus,
		Occupation:   &input.Occupation,
		EmployerName: &input.EmployerName,
		WorkAddress:  &input.WorkAddress,
		Type:         &input.EmploymentType,
	})
	if err != nil {
		return err
	}
	fmt.Println("Employment created")

	// 6. Create emergency
	_, err = qtx.CreateEmergencyContact(ctx, repository.CreateEmergencyContactParams{
		ResidentID: residentID,
		Name:       input.EmergencyName,
		Relation:   input.EmergencyRelation,
		Phone:      input.EmergencyPhone,
		Email:      &input.EmergencyEmail,
	})
	if err != nil {
		return err
	}
	fmt.Println("Emergency created")

	// 6. Create additional
	_, err = qtx.CreateAdditional(ctx, repository.CreateAdditionalParams{
		ResidentID:      residentID,
		Religion:        &input.Religion,
		Ethnicity:       &input.Ethnicity,
		NationalID:      &input.NationalID,
		Disability:      &input.Disability,
		EducationLevel:  &input.EducationLevel,
		MaritalStatus:   &input.MaritalStatus,
		LanguagesSpoken: input.LanguagesSpoken,
	})
	if err != nil {
		return err
	}
	fmt.Println("Additional created")

	// 7. Create payment
	_, err = qtx.CreatePayment(ctx, repository.CreatePaymentParams{
		ResidentID: residentID,
		Status:     "unpaid",
	})
	if err != nil {
		return err
	}
	fmt.Println("Payment created")

	return tx.Commit(ctx)
}

func (s *Service) GetResident(ctx context.Context, id uuid.UUID) (repository.GetResidentRow, error) {
	return s.q.GetResident(ctx, id)
}

func (s *Service) GetResidentDocuments(ctx context.Context, id uuid.UUID) ([]repository.Document, error) {
	return s.q.GetResidentDocuments(ctx, id)
}

func (s *Service) GetResidentPayment(ctx context.Context, id uuid.UUID) (repository.Payment, error) {
	return s.q.GetPaymentByResident(ctx, id)
}

func (s *Service) GetResidents(ctx context.Context, limit, offset int) (int64, []repository.ListResidentsRow, error) {
	count, err := s.q.CountListResidents(ctx)
	if err != nil {
		return 0, nil, err
	}

	cities, err := s.q.ListResidents(ctx, repository.ListResidentsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return 0, nil, err
	}

	return count, cities, nil
}

func (s *Service) SearchResidents(ctx context.Context, limit, offset int, query string) (int64, []repository.SearchResidentsRow, error) {
	count, err := s.q.CountSearchResidents(ctx, query)
	if err != nil {
		return 0, nil, err
	}

	cities, err := s.q.SearchResidents(ctx, repository.SearchResidentsParams{
		Query:  query,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return 0, nil, err
	}

	return count, cities, nil
}

func (s *Service) GetUnpaidResidents(ctx context.Context, limit, offset int) (int64, []repository.ListUnpaidResidentsRow, error) {
	count, err := s.q.CountListUnpaidResidents(ctx)
	if err != nil {
		return 0, nil, err
	}

	cities, err := s.q.ListUnpaidResidents(ctx, repository.ListUnpaidResidentsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return 0, nil, err
	}

	return count, cities, nil
}

func (s *Service) SearchUnpaidResidents(ctx context.Context, limit, offset int, query string) (int64, []repository.SearchUnpaidResidentsRow, error) {
	count, err := s.q.CountSearchUnpaidResidents(ctx, query)
	if err != nil {
		return 0, nil, err
	}

	cities, err := s.q.SearchUnpaidResidents(ctx, repository.SearchUnpaidResidentsParams{
		Query:  query,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return 0, nil, err
	}

	return count, cities, nil
}

func (s *Service) GetUnverifiedResidents(ctx context.Context, limit, offset int) (int64, []repository.ListUnverifiedResidentsRow, error) {
	count, err := s.q.CountListUnverifiedResidents(ctx)
	if err != nil {
		return 0, nil, err
	}

	cities, err := s.q.ListUnverifiedResidents(ctx, repository.ListUnverifiedResidentsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return 0, nil, err
	}

	return count, cities, nil
}

func (s *Service) SearchUnverifiedResidents(ctx context.Context, limit, offset int, query string) (int64, []repository.SearchUnverifiedResidentsRow, error) {
	count, err := s.q.CountSearchUnverifiedResidents(ctx, query)
	if err != nil {
		return 0, nil, err
	}

	cities, err := s.q.SearchUnverifiedResidents(ctx, repository.SearchUnverifiedResidentsParams{
		Query:  query,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return 0, nil, err
	}

	return count, cities, nil
}

func (s *Service) UpdatePaymentInfo(ctx context.Context, id uuid.UUID, amount float64, status, method, description string, receipt *string) error {
	_, err := s.q.UpdatePayment(ctx, repository.UpdatePaymentParams{
		ID:          id,
		Amount:      &amount,
		Status:      &status,
		Method:      &method,
		Description: &description,
		Reference:   receipt,
	})
	return err
}

func (s *Service) UpdateDocumentInfo(ctx context.Context, id uuid.UUID, status string, url *string) error {
	_, err := s.q.UpdateDocument(ctx, repository.UpdateDocumentParams{
		ID:     id,
		Status: &status,
		Url:    url,
	})
	return err
}

func (s *Service) UpdatePersonalInfo(ctx context.Context, id uuid.UUID, input types.ResidentPersonalPayload) error {
	return s.q.UpdateResident(ctx, repository.UpdateResidentParams{
		ID:         id,
		Email:      input.Email,
		FirstName:  input.FirstName,
		SecondName: input.SecondName,
		LastName:   input.LastName,
		BirthDate:  input.BirthDate,
		Gender:     input.Gender,
		Phone:      input.Phone,
	})
}

func (s *Service) UpdateAdditionalInfo(ctx context.Context, id uuid.UUID, input types.ResidentAdditionalPayload) error {
	_, err := s.q.UpdateAdditional(ctx, repository.UpdateAdditionalParams{
		ID:              id,
		Religion:        &input.Religion,
		Ethnicity:       &input.Ethnicity,
		NationalID:      &input.NationalID,
		Disability:      &input.Disability,
		EducationLevel:  &input.EducationLevel,
		LanguagesSpoken: input.LanguagesSpoken,
		MaritalStatus:   &input.MaritalStatus,
	})
	return err
}

func (s *Service) UpdateEmploymentInfo(ctx context.Context, id uuid.UUID, input types.ResidentEmploymentPayload) error {
	_, err := s.q.UpdateEmployment(ctx, repository.UpdateEmploymentParams{
		ID:           id,
		Status:       &input.EmploymentStatus,
		Type:         &input.EmploymentType,
		Occupation:   &input.Occupation,
		EmployerName: &input.EmployerName,
		WorkAddress:  &input.WorkAddress,
	})
	return err
}

func (s *Service) UpdateEmergencyContact(ctx context.Context, id uuid.UUID, input types.ResidentEmergencyPayload) error {
	_, err := s.q.UpdateEmergencyContact(ctx, repository.UpdateEmergencyContactParams{
		ID:       id,
		Name:     &input.EmergencyName,
		Relation: &input.EmergencyRelation,
		Phone:    &input.EmergencyPhone,
		Email:    &input.EmergencyEmail,
	})
	return err
}

func (s *Service) GetAddressInfo(ctx context.Context, id uuid.UUID) (repository.GetAddressRow, error) {
	return s.q.GetAddress(ctx, id)
}

func (s *Service) UpdateAddressInfo(ctx context.Context, id uuid.UUID, houseNumber string, kebeleID, subcityID, cityID uuid.UUID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	addr, err := qtx.GetAddressByLocations(ctx, repository.GetAddressByLocationsParams{
		HouseNumber: houseNumber,
		KebeleID:    kebeleID,
		SubcityID:   subcityID,
		CityID:      cityID,
	})
	if err != nil {
		fmt.Println("Address not found")
		newAddr, err := qtx.CreateAddress(ctx, repository.CreateAddressParams{
			HouseNumber: houseNumber,
			KebeleID:    kebeleID,
			SubcityID:   subcityID,
			CityID:      cityID,
		})
		if err != nil {
			return err
		}
		fmt.Println("New address created")
		err = qtx.UpdateResidentAddress(ctx, repository.UpdateResidentAddressParams{
			AddressID: &(newAddr.ID),
			ID:        id,
		})
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Address found")
		err = qtx.UpdateResidentAddress(ctx, repository.UpdateResidentAddressParams{
			AddressID: &addr.ID,
			ID:        id,
		})
		if err != nil {
			return err
		}
	}
	fmt.Println("Committing transaction")
	return tx.Commit(ctx)
}

func (s *Service) ReplaceDocuments(ctx context.Context, id uuid.UUID, docsUrl []string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	err = qtx.DeleteResidentDocuments(ctx, id)
	if err != nil {
		return err
	}

	for _, doc := range docsUrl {
		_, err = qtx.CreateDocument(ctx, repository.CreateDocumentParams{
			Url:    doc,
			Status: "Pending",
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *Service) UpdateBiometricInfo(ctx context.Context, id uuid.UUID, bloodType string, faceUrl *string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	err = qtx.UpdateBiometric(ctx, repository.UpdateBiometricParams{
		BloodType:  &bloodType,
		FaceUrl:    faceUrl,
		ResidentID: id,
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
