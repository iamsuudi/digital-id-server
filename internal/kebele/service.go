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
		Lat:       input.Lat,
		Lon:       input.Lon,
	})
}

func (s *Service) DeleteKebele(ctx context.Context, id uuid.UUID) error {
	return s.q.SoftDeleteKebele(ctx, id)
}

func (s *Service) RemoveStaff(ctx context.Context, staffID uuid.UUID) error {
	return s.q.RevokeUserPlacement(ctx, repository.RevokeUserPlacementParams{
		KebeleID: nil,
		ID:       staffID,
	})
}

func (s *Service) AddStaff(ctx context.Context, kebeleID, staffID uuid.UUID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	kebele, err := qtx.GetKebele(ctx, kebeleID)
	if err != nil {
		return err
	}

	err = qtx.GrantUserPlacement(ctx, repository.GrantUserPlacementParams{
		ID:        staffID,
		CityID:    &(kebele.CityID),
		SubcityID: &(kebele.SubcityID),
		KebeleID:  &(kebeleID),
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Service) AssignExecutive(ctx context.Context, kebeleID, executiveID uuid.UUID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	kebele, err := qtx.GetKebele(ctx, kebeleID)
	if err != nil {
		return err
	}

	// Remove previous executive
	if kebele.ExecutiveID != nil {
		err = qtx.RevokeUserPlacement(ctx, repository.RevokeUserPlacementParams{
			KebeleID: nil,
			ID:       *kebele.ExecutiveID,
		})
		if err != nil {
			return err
		}
	}

	err = qtx.GrantUserPlacement(ctx, repository.GrantUserPlacementParams{
		ID:        executiveID,
		CityID:    &(kebele.CityID),
		SubcityID: &(kebele.SubcityID),
		KebeleID:  &(kebeleID),
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Service) UpdateKebele(ctx context.Context, id uuid.UUID, name string, lat, lon *float64) error {
	_, err := s.q.UpdateKebele(ctx, repository.UpdateKebeleParams{
		ID:   id,
		Name: name,
		Lat:  lat,
		Lon:  lon,
	})
	return err
}

func (s *Service) GetKebele(ctx context.Context, id uuid.UUID) (repository.GetKebeleRow, error) {
	return s.q.GetKebele(ctx, id)
}

func (s *Service) GetEncoders(ctx context.Context, id uuid.UUID) ([]repository.ListUsersByKebeleAndRoleRow, error) {
	return s.q.ListUsersByKebeleAndRole(ctx, repository.ListUsersByKebeleAndRoleParams{
		KebeleID: &id,
		RoleSlug: "encoder",
	})
}

func (s *Service) GetCashiers(ctx context.Context, id uuid.UUID) ([]repository.ListUsersByKebeleAndRoleRow, error) {
	return s.q.ListUsersByKebeleAndRole(ctx, repository.ListUsersByKebeleAndRoleParams{
		KebeleID: &id,
		RoleSlug: "cashier",
	})
}

func (s *Service) GetKebeles(ctx context.Context, limit, offset int) (int64, []repository.ListKebelesRow, error) {
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
