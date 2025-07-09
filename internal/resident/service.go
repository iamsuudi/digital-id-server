package resident

import (
	"context"
	"errors"

	"github.com/iamsuudi/digital-id-server/database"
	"github.com/iamsuudi/digital-id-server/database/sqlc"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	q  *sqlc.Queries
	db sqlc.DBTX // interface allowing transaction
}

func NewService(q *sqlc.Queries, db sqlc.DBTX) *Service {
	return &Service{q: q, db: db}
}

func (s *Service) RegisterResident(ctx context.Context, input RegisterResidentInput, faceURL, docURL string) error {
	return database.WithTx(ctx, s.db, func(tx *sqlc.Queries) error {
		// 1. Create resident
		residentID, err := tx.CreateResident(ctx, sqlc.CreateResidentParams{
			Email:            input.Email,
			FirstName:        input.FirstName,
			SecondName:       input.SecondName,
			LastName:         input.LastName,
			BirthDate:        toPgTimestamp(input.BirthDate),
			Gender:           toGender(input.Gender),
			Phone:            input.Phone,
			MaritalStatus:    toMarital(input.MaritalStatus),
			Religion:         toReligion(input.Religion),
			Ethnicity:        toPgText(input.Ethnicity),
			DisabilityStatus: toPgText(input.DisabilityStatus),
			EducationLevel:   toPgText(input.EducationLevel),
			LanguagesSpoken:  input.LanguagesSpoken,
		})
		if err != nil {
			return err
		}

		// 2. Create biometric
		err = tx.CreateBiometric(ctx, sqlc.CreateBiometricParams{
			ResidentID:  residentID,
			BloodType:   input.BloodType,
			Face:        faceURL,
			Fingerprint: nil, // not used here
		})
		if err != nil {
			return err
		}

		// 3. Get or create address
		addr, err := tx.GetAddress(ctx, sqlc.GetAddressParams{
			HouseNumber: input.HouseNumber,
			District:    input.District,
			CityID:      input.CityID,
		})
		if errors.Is(err, pgx.ErrNoRows) {
			newAddrID, err := tx.CreateAddress(ctx, sqlc.CreateAddressParams{
				HouseNumber: input.HouseNumber,
				District:    input.District,
				CityID:      input.CityID,
			})
			if err != nil {
				return err
			}
			err = tx.UpdateResidentAddress(ctx, sqlc.UpdateResidentAddressParams{
				AddressID: newAddrID,
				ID:        residentID,
			})
			if err != nil {
				return err
			}
		} else if err == nil {
			err = tx.UpdateResidentAddress(ctx, sqlc.UpdateResidentAddressParams{
				AddressID: addr.ID,
				ID:        residentID,
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}

		// 4. Create document
		err = tx.CreateDocument(ctx, sqlc.CreateDocumentParams{
			Type:       toDocumentType(input.DocumentType),
			Number:     input.DocumentNumber,
			ResidentID: residentID,
			Url:        docURL,
			Status:     "pending",
		})
		if err != nil {
			return err
		}

		// 5. Create employment
		err = tx.CreateEmployment(ctx, sqlc.CreateEmploymentParams{
			ResidentID:   residentID,
			Status:       input.EmploymentStatus,
			Occupation:   toPgText(input.Occupation),
			EmployerName: toPgText(input.EmployerName),
			WorkAddress:  toPgText(input.WorkAddress),
		})
		if err != nil {
			return err
		}

		// 6. Create emergency
		err = tx.CreateEmergency(ctx, sqlc.CreateEmergencyParams{
			ResidentID: residentID,
			Name:       input.EmergencyName,
			Relation:   input.EmergencyRelation,
			Phone:      input.EmergencyPhone,
		})
		return err
	})
}

func (s *Service) GetResident(ctx context.Context, id int32) (sqlc.GetResidentFullRow, error) {
	return s.q.GetResidentFull(ctx, id)
}
