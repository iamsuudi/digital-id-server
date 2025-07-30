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
	return s.q.CreateCity(ctx, input.Name)
}

func (s *Service) UpdateCity(ctx context.Context, id uuid.UUID, input types.CityInput) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	_, err = qtx.UpdateCity(ctx, repository.UpdateCityParams{
		ID:   id,
		Name: input.Name,
	})
	if err != nil {
		return err
	}

	err = qtx.RevokeUserPlacement(ctx, repository.RevokeUserPlacementParams{
		CityID:    &id,
		SubcityID: nil,
		KebeleID:  nil,
		RoleSlug:  "admin",
	})
	if err != nil {
		return err
	}

	if input.AdminID != nil {

		if input.AdminID != nil {
			err = qtx.GrantUserPlacement(ctx, repository.GrantUserPlacementParams{
				ID:        *input.AdminID,
				CityID:    &id,
				SubcityID: nil,
				KebeleID:  nil,
			})
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
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
