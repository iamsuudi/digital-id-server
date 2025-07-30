package kebele

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

func (s *Service) CreateKebele(ctx context.Context, input types.KebeleInput) (repository.Kebele, error) {
	return s.q.CreateKebele(ctx, repository.CreateKebeleParams{
		Name:      input.Name,
		CityID:    input.CityID,
		SubcityID: input.SubCityID,
	})
}

func (s *Service) UpdateKebele(ctx context.Context, id uuid.UUID, input types.KebeleInput) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	_, err = qtx.UpdateKebele(ctx, repository.UpdateKebeleParams{
		ID:        id,
		Name:      input.Name,
		CityID:    input.CityID,
		SubcityID: input.SubCityID,
	})
	if err != nil {
		return err
	}
	
	err = qtx.RevokeUserPlacement(ctx, repository.RevokeUserPlacementParams{
		CityID:    &input.CityID,
		SubcityID: input.SubCityID,
		KebeleID:  &id,
		RoleSlug: "executive",
	})
	if err != nil {
		return err
	}

	if input.ExecutiveID != nil {
		err = qtx.GrantUserPlacement(ctx, repository.GrantUserPlacementParams{
			ID:        *input.ExecutiveID,
			CityID:    &input.CityID,
			SubcityID: input.SubCityID,
			KebeleID:  &id,
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *Service) GetKebele(ctx context.Context, id uuid.UUID) (repository.GetKebeleRow, error) {
	return s.q.GetKebele(ctx, id)
}

func (s *Service) GetAll(ctx context.Context, limit, offset int) (int64, []repository.ListKebelesRow, error) {
	count, err := s.q.CountListSubcities(ctx)
	if err != nil {
		return 0, nil, err
	}

	cities, err := s.q.ListKebeles(ctx, repository.ListKebelesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return 0, nil, err
	}

	return count, cities, nil
}

func (s *Service) SearchKebeles(ctx context.Context, limit, offset int, query string) (int64, []repository.SearchKebelesRow, error) {
	count, err := s.q.CountSearchKebeles(ctx, query)
	if err != nil {
		return 0, nil, err
	}

	kebeles, err := s.q.SearchKebeles(ctx, repository.SearchKebelesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
		Query:  query,
	})
	if err != nil {
		return 0, nil, err
	}

	return count, kebeles, nil
}
