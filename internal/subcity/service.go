package subcity

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

func (s *Service) CreateSubCity(ctx context.Context, input types.SubCityInput) (repository.Subcity, error) {
	return s.q.CreateSubCity(ctx, repository.CreateSubCityParams{
		Name:   input.Name,
		CityID: input.CityID,
	})
}

func (s *Service) DeleteSubCity(ctx context.Context, id uuid.UUID) error {
	return s.q.SoftDeleteSubCity(ctx, id)
}

func (s *Service) RemoveStaff(ctx context.Context, staffID uuid.UUID) error {
	return s.q.RevokeUserPlacement(ctx, repository.RevokeUserPlacementParams{
		SubcityID: nil,
		ID:        staffID,
	})
}

func (s *Service) AddStaff(ctx context.Context, subCityID, staffID uuid.UUID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	subcity, err := qtx.GetSubCity(ctx, subCityID)
	if err != nil {
		return err
	}

	err = qtx.GrantUserPlacement(ctx, repository.GrantUserPlacementParams{
		ID:        staffID,
		CityID:    &(subcity.CityID),
		SubcityID: &subCityID,
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Service) AssignManager(ctx context.Context, subCityID, managerID uuid.UUID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	subcity, err := qtx.GetSubCity(ctx, subCityID)
	if err != nil {
		return err
	}

	// Remove previous manager
	if subcity.ManagerID != nil {
		err = qtx.RevokeUserPlacement(ctx, repository.RevokeUserPlacementParams{
			SubcityID: nil,
			ID:        *subcity.ManagerID,
		})
		if err != nil {
			return err
		}
	}

	err = qtx.GrantUserPlacement(ctx, repository.GrantUserPlacementParams{
		ID:        managerID,
		CityID:    &(subcity.CityID),
		SubcityID: &subCityID,
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Service) UpdateSubCity(ctx context.Context, id uuid.UUID, name string, lat, lon *float64) error {
	_, err := s.q.UpdateSubCity(ctx, repository.UpdateSubCityParams{
		ID:   id,
		Name: name,
		Lat:  lat,
		Lon:  lon,
	})
	return err
}

func (s *Service) GetSubCity(ctx context.Context, id uuid.UUID) (repository.GetSubCityRow, error) {
	return s.q.GetSubCity(ctx, id)
}

func (s *Service) GetKebeles(ctx context.Context, id uuid.UUID) ([]repository.GetKebelesForSubCityRow, error) {
	return s.q.GetKebelesForSubCity(ctx, id)
}

func (s *Service) GetSubCities(ctx context.Context, limit, offset int) (int64, []repository.ListSubCitiesRow, error) {
	count, err := s.q.CountListSubcities(ctx)
	if err != nil {
		return 0, nil, err
	}

	cities, err := s.q.ListSubCities(ctx, repository.ListSubCitiesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return 0, nil, err
	}

	return count, cities, nil
}

func (s *Service) SearchSubCities(ctx context.Context, limit, offset int, query string) (int64, []repository.SearchSubCitiesRow, error) {
	count, err := s.q.CountSearchSubCities(ctx, query)
	if err != nil {
		return 0, nil, err
	}

	subcities, err := s.q.SearchSubCities(ctx, repository.SearchSubCitiesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
		Query:  query,
	})
	if err != nil {
		return 0, nil, err
	}

	return count, subcities, nil
}
