package city

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

func (s *Service) CreateCity(ctx context.Context, input types.CityInput) (repository.City, error) {
	return s.q.CreateCity(ctx, repository.CreateCityParams{
		Name: input.Name,
		Lat: &input.Lat,
		Lon: &input.Lon,
	})
}

func (s *Service) DeleteCity(ctx context.Context, id uuid.UUID) error {
	return s.q.SoftDeleteCity(ctx, id)
}

func (s *Service) RemoveStaff(ctx context.Context, staffID uuid.UUID) error {
	return s.q.RevokeUserPlacement(ctx, repository.RevokeUserPlacementParams{
		CityID: nil,
		ID:     staffID,
	})
}

func (s *Service) AddStaff(ctx context.Context, CityID, staffID uuid.UUID) error {
	return s.q.GrantUserPlacement(ctx, repository.GrantUserPlacementParams{
		ID:     staffID,
		CityID: &(CityID),
	})
}

func (s *Service) AssignAdmin(ctx context.Context, cityID, adminID uuid.UUID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	city, err := qtx.GetCity(ctx, cityID)
	if err != nil {
		return err
	}

	// Remove previous admin
	if city.AdminID != nil {
		err = qtx.RevokeUserPlacement(ctx, repository.RevokeUserPlacementParams{
			CityID: nil,
			ID:     *city.AdminID,
		})
		if err != nil {
			return err
		}
	}

	err = qtx.GrantUserPlacement(ctx, repository.GrantUserPlacementParams{
		ID:     adminID,
		CityID: &cityID,
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Service) UpdateCityInfo(ctx context.Context, id uuid.UUID, name string, lat, lon *float64) error {
	return s.q.UpdateCity(ctx, repository.UpdateCityParams{
		ID:   id,
		Name: name,
		Lat:  lat,
		Lon:  lon,
	})
}

func (s *Service) GetSubCities(ctx context.Context, id uuid.UUID) ([]repository.GetSubCitiesForCityRow, error) {
	return s.q.GetSubCitiesForCity(ctx, id)
}

func (s *Service) GetCity(ctx context.Context, id uuid.UUID) (repository.GetCityRow, error) {
	return s.q.GetCity(ctx, id)
}

func (s *Service) GetCities(ctx context.Context, limit, offset int) (int64, []repository.ListCitiesRow, error) {
	count, err := s.q.CountListCities(ctx)
	if err != nil {
		fmt.Println(err)
		return 0, nil, err
	}

	cities, err := s.q.ListCities(ctx, repository.ListCitiesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		fmt.Println(err)
		return 0, nil, err
	}

	return count, cities, nil
}

func (s *Service) SearchCities(ctx context.Context, limit, offset int, query string) (int64, []repository.SearchCitiesRow, error) {
	count, err := s.q.CountSearchCities(ctx, query)
	if err != nil {
		return 0, nil, err
	}

	cities, err := s.q.SearchCities(ctx, repository.SearchCitiesParams{
		Query:  query,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return 0, nil, err
	}

	return count, cities, nil
}
