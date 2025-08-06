package resident

import (
	"context"

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

func (s *Service) RegisterResident(ctx context.Context, input types.ResidentInput, faceURL, docURL string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	// 1. Create resident
	_, err = qtx.CreateResident(ctx, repository.CreateResidentParams{
		Email:           input.Email,
		FirstName:       input.FirstName,
		SecondName:      input.SecondName,
		LastName:        input.LastName,
		BirthDate:       input.BirthDate,
		Gender:          input.Gender,
		Phone:           input.Phone,
		MaritalStatus:   &input.MaritalStatus,
		Religion:        &input.Religion,
		Ethnicity:       &input.Ethnicity,
		Disability:      &input.Disability,
		EducationLevel:  &input.EducationLevel,
		LanguagesSpoken: input.LanguagesSpoken,
	})
	if err != nil {
		return err
	}

	// // 2. Create biometric
	// err = qtx.CreateBiometric(ctx, repository.CreateBiometricParams{
	// 	ResidentID: residentID,
	// 	BloodType:  input.BloodType,
	// 	Face:       faceURL,
	// })
	// if err != nil {
	// 	return err
	// }

	// // 3. Get or create address
	// addr, err := qtx.GetAddress(ctx, repository.GetAddressParams{
	// 	HouseNumber: input.HouseNumber,
	// 	KebeleID:    input.KebeleID,
	// 	CityID:      input.CityID,
	// })
	// if errors.Is(err, pgx.ErrNoRows) {
	// 	newAddrID, err := qtx.CreateAddress(ctx, repository.CreateAddressParams{
	// 		HouseNumber: input.HouseNumber,
	// 		KebeleID:    input.KebeleID,
	// 		CityID:      input.CityID,
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}
	// 	err = qtx.UpdateResidentAddress(ctx, repository.UpdateResidentAddressParams{
	// 		AddressID: &newAddrID,
	// 		ID:        residentID,
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}
	// } else if err == nil {
	// 	err = qtx.UpdateResidentAddress(ctx, repository.UpdateResidentAddressParams{
	// 		AddressID: &addr.ID,
	// 		ID:        residentID,
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}
	// } else {
	// 	return err
	// }

	// // 4. Create document
	// err = qtx.CreateDocument(ctx, repository.CreateDocumentParams{
	// 	Type:       input.DocumentType,
	// 	Number:     input.DocumentNumber,
	// 	ResidentID: residentID,
	// 	Url:        docURL,
	// 	Status:     "pending",
	// })
	// if err != nil {
	// 	return err
	// }

	// // 5. Create employment
	// err = qtx.CreateEmployment(ctx, repository.CreateEmploymentParams{
	// 	ResidentID:   residentID,
	// 	Status:       input.EmploymentStatus,
	// 	Occupation:   &input.Occupation,
	// 	EmployerName: &input.EmployerName,
	// 	WorkAddress:  &input.WorkAddress,
	// })
	// if err != nil {
	// 	return err
	// }

	// // 6. Create emergency
	// err = qtx.CreateEmergency(ctx, repository.CreateEmergencyParams{
	// 	ResidentID: residentID,
	// 	Name:       input.EmergencyName,
	// 	Relation:   input.EmergencyRelation,
	// 	Phone:      input.EmergencyPhone,
	// })
	// if err != nil {
	// 	return err
	// }

	return tx.Commit(ctx)
}

func (s *Service) GetResident(ctx context.Context, id uuid.UUID) (repository.GetResidentRow, error) {
	return s.q.GetResident(ctx, id)
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
