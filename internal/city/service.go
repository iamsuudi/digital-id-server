package city

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"digital-id-server/shared/types"
	"digital-id-server/internal/repository"
)

type Service struct {
	db *pgxpool.Pool
	q  *repository.Queries
}

func NewService(dbConn *pgxpool.Pool, dbQueries *repository.Queries) *Service {
	return &Service{db: dbConn, q: dbQueries}
}

func (s *Service) CreateCity(ctx context.Context, input types.CityInput) (repository.City, error) {
	return s.q.CreateCity(ctx, input.Name)
}

func (s *Service) UpdateCity(ctx context.Context, id uuid.UUID, input types.CityInput) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	if input.AdminID != nil {
		err := qtx.UpdateUserPlacement(ctx, repository.UpdateUserPlacementParams{
			ID:        *input.AdminID,
			CityID:    &id,
			SubcityID: nil,
			KebeleID:  nil,
		})
		if err != nil {
			return err
		}
	}

	_, err = qtx.UpdateCity(ctx, repository.UpdateCityParams{
		ID:     id,
		Name:   input.Name,
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Service) GetCity(ctx context.Context, id uuid.UUID) (repository.GetCityRow, error) {
	return s.q.GetCity(ctx, id)
}

func (s *Service) GetAll(ctx context.Context, limit, offset int) (int64, []repository.ListCitiesRow, error) {
	count, err := s.q.CountListCities(ctx)
	if err != nil {
		return 0, nil, err
	}

	cities, err := s.q.ListCities(ctx, repository.ListCitiesParams{
		Limit: int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return 0, nil, err
	}

	return count, cities, nil
}
