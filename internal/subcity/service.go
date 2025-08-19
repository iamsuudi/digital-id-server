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

func (s *Service) CreateSubCity(ctx context.Context, actorID uuid.UUID, input types.SubCityInput) (repository.Subcity, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repository.Subcity{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	subcity, err := qtx.CreateSubCity(ctx, repository.CreateSubCityParams{
		Name:   input.Name,
		CityID: input.CityID,
	})
	if err != nil {
		return repository.Subcity{}, err
	}

	qtx.InsertAuditLog(ctx, repository.InsertAuditLogParams{
		ActorID:         actorID,
		ActionType:      "CREATE_SUBCITY",
		TargetSubcityID: &(subcity.ID),
		ObjectType:      "subcity",
		Diff: map[string]any{
			"after": subcity,
		},
	})

	return subcity, tx.Commit(ctx)
}

func (s *Service) DeleteSubCity(ctx context.Context, actorID, id uuid.UUID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := s.q.WithTx(tx)

	// 1. Before
	before, err := qtx.GetSubCity(ctx, id)
	if err != nil {
		return err
	}

	// 2. After
	err = qtx.SoftDeleteSubCity(ctx, id)
	if err != nil {
		return err
	}

	// 3. Insert audit log
	err = qtx.InsertAuditLog(ctx, repository.InsertAuditLogParams{
		ActorID:    actorID,
		ActionType: "DELETE_SUBCITY",
		ObjectType: "subcity",
		Diff: map[string]any{
			"before": before,
		},
	})

	return tx.Commit(ctx)
}

func (s *Service) RemoveStaff(ctx context.Context, actorID, staffID uuid.UUID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := s.q.WithTx(tx)

	//  Before
	user, err := qtx.GetUserByID(ctx, staffID)
	if err != nil {
		return err
	}
	subcity, err := qtx.GetSubCity(ctx, *user.SubcityID)
	if err != nil {
		return err
	}

	// After
	err = qtx.RevokeUserPlacement(ctx, repository.RevokeUserPlacementParams{
		SubcityID: nil,
		ID:        staffID,
	})

	// Insert audit log
	err = qtx.InsertAuditLog(ctx, repository.InsertAuditLogParams{
		ActorID:         actorID,
		ActionType:      "REMOVE_MANAGER",
		TargetSubcityID: &(subcity.ID),
		ObjectType:      "subcity",
		Diff: map[string]any{
			"before": map[string]any{
				"user_name": subcity.ManagerName,
				"subcity":   subcity.Name,
				"position":  "Manager",
			},
		},
	})

	return tx.Commit(ctx)
}

func (s *Service) AssignManager(ctx context.Context, actorID, subCityID, managerID uuid.UUID) error {
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

	// 1. Remove previous manager
	if subcity.ManagerID != nil {
		err = qtx.RevokeUserPlacement(ctx, repository.RevokeUserPlacementParams{
			SubcityID: nil,
			ID:        *subcity.ManagerID,
		})
		if err != nil {
			return err
		}
	}

	// 2. Assign new manager
	err = qtx.GrantUserPlacement(ctx, repository.GrantUserPlacementParams{
		ID:        managerID,
		CityID:    &(subcity.CityID),
		SubcityID: &subCityID,
	})
	if err != nil {
		return err
	}

	user, err := qtx.GetUserByID(ctx, managerID)
	if err != nil {
		return err
	}

	// 3. Insert audit log
	err = qtx.InsertAuditLog(ctx, repository.InsertAuditLogParams{
		ActorID:         actorID,
		ActionType:      "ASSIGN_MANAGER",
		TargetSubcityID: &(subCityID),
		ObjectType:      "subcity",
		Diff: map[string]any{
			"before": map[string]any{
				"user_name": subcity.ManagerName,
				"subcity":   subcity.Name,
				"position":  "Manager",
			},
			"after": map[string]any{
				"user_name": user.FullName,
				"subcity":   subcity.Name,
				"position":  "Manager",
			},
		},
	})

	return tx.Commit(ctx)
}

func (s *Service) UpdateSubCity(ctx context.Context, actorID, id uuid.UUID, name string, lat, lon *float64) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := s.q.WithTx(tx)

	// 1. Before
	before, err := qtx.GetSubCity(ctx, id)
	if err != nil {
		return err
	}

	// 2. After
	updated, err := qtx.UpdateSubCity(ctx, repository.UpdateSubCityParams{
		ID:   id,
		Name: name,
		Lat:  lat,
		Lon:  lon,
	})

	// 3. Insert audit log
	err = qtx.InsertAuditLog(ctx, repository.InsertAuditLogParams{
		ActorID:         actorID,
		TargetSubcityID: &id,
		ActionType:      "UPDATE_SUBCITY",
		ObjectType:      "subcity",
		Diff: map[string]any{
			"before": before,
			"after":  updated,
		},
	})
	return tx.Commit(ctx)
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
