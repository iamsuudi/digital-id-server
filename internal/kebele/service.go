package kebele

import (
	"context"
	"fmt"
	"strings"

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

func (s *Service) CreateKebele(ctx context.Context, actorID uuid.UUID, input types.KebeleInput) (repository.Kebele, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repository.Kebele{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	kebele, err := qtx.CreateKebele(ctx, repository.CreateKebeleParams{
		Name:      input.Name,
		Lat:       input.Lat,
		Lon:       input.Lon,
		SubcityID: input.SubCityID,
		CityID:    input.CityID,
	})
	if err != nil {
		return repository.Kebele{}, err
	}

	err = qtx.InsertAuditLog(ctx, repository.InsertAuditLogParams{
		ActorID:        actorID,
		ActionType:     "CREATE_KEBELE",
		TargetKebeleID: &(kebele.ID),
		ObjectType:     "kebele",
		Diff: map[string]any{
			"after": kebele,
		},
	})
	if err != nil {
		return repository.Kebele{}, err
	}

	return kebele, tx.Commit(ctx)
}

func (s *Service) DeleteKebele(ctx context.Context, actorID, id uuid.UUID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	// 1. Before
	kebele, err := qtx.GetKebele(ctx, id)
	if err != nil {
		return err
	}

	// 2. Delete
	err = s.q.SoftDeleteKebele(ctx, id)
	if err != nil {
		return err
	}

	// 3. Insert audit log
	qtx.InsertAuditLog(ctx, repository.InsertAuditLogParams{
		ActorID:    actorID,
		ActionType: "DELETE_KEBELE",
		ObjectType: "kebele",
		Diff: map[string]any{
			"before": kebele,
		},
	})
	return tx.Commit(ctx)
}

func (s *Service) RemoveStaff(ctx context.Context, actorID, kebeleID, staffID uuid.UUID) error {
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
	kebele, err := qtx.GetKebele(ctx, kebeleID)
	if err != nil {
		return err
	}

	// After
	err = qtx.RevokeUserPlacement(ctx, repository.RevokeUserPlacementParams{
		KebeleID: nil,
		ID:       staffID,
	})

	// Insert audit log
	err = qtx.InsertAuditLog(ctx, repository.InsertAuditLogParams{
		ActorID:        actorID,
		ActionType:     fmt.Sprintf("REMOVE_%s", strings.ToUpper(user.RoleSlug)),
		TargetKebeleID: &(kebele.ID),
		ObjectType:     "kebele",
		Diff: map[string]any{
			"before": map[string]any{
				"user_name": user.FullName,
				"kebele":    kebele.Name,
				"position":  user.RoleSlug,
			},
		},
	})

	return tx.Commit(ctx)
}

func (s *Service) AddStaff(ctx context.Context, actorID, kebeleID, staffID uuid.UUID) error {
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

	user, err := qtx.GetUserByID(ctx, staffID)
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

	err = qtx.InsertAuditLog(ctx, repository.InsertAuditLogParams{
		ActorID:        actorID,
		ActionType:     fmt.Sprintf("ADD_%s", strings.ToUpper(user.RoleSlug)),
		TargetKebeleID: &(kebele.ID),
		ObjectType:     "kebele",
		Diff: map[string]any{
			"after": map[string]any{
				"user_name": user.FullName,
				"kebele":    kebele.Name,
				"position":  user.RoleSlug,
			},
		},
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Service) AssignExecutive(ctx context.Context, actorID, kebeleID, executiveID uuid.UUID) error {
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

	user, err := qtx.GetUserByID(ctx, executiveID)
	if err != nil {
		return err
	}

	// Assign new executive
	err = qtx.GrantUserPlacement(ctx, repository.GrantUserPlacementParams{
		ID:        executiveID,
		CityID:    &(kebele.CityID),
		SubcityID: &(kebele.SubcityID),
		KebeleID:  &(kebeleID),
	})
	if err != nil {
		return err
	}

	err = qtx.InsertAuditLog(ctx, repository.InsertAuditLogParams{
		ActorID:        actorID,
		ActionType:     "ASSIGN_EXECUTIVE",
		ObjectType:     "kebele",
		TargetKebeleID: &(kebeleID),
		Diff: map[string]any{
			"after": map[string]any{
				"user_name": user.FullName,
				"kebele":    kebele.Name,
				"position":  "Executive",
			},
		},
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Service) UpdateKebeleInfo(ctx context.Context, actorID, id uuid.UUID, name string, lat, lon *float64) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := s.q.WithTx(tx)

	// 1. Before
	before, err := qtx.GetKebele(ctx, id)
	if err != nil {
		return err
	}

	// 2. After
	updated, err := qtx.UpdateKebele(ctx, repository.UpdateKebeleParams{
		ID:   id,
		Name: name,
		Lat:  lat,
		Lon:  lon,
	})
	if err != nil {
		return err
	}

	// 3. Insert audit log
	err = qtx.InsertAuditLog(ctx, repository.InsertAuditLogParams{
		ActorID:        actorID,
		TargetKebeleID: &id,
		ActionType:     "UPDATE_KEBELE",
		ObjectType:     "kebele",
		Diff: map[string]any{
			"before": before,
			"after":  updated,
		},
	})
	return tx.Commit(ctx)
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

func (s *Service) ListAuditLogs(ctx context.Context, limit, offset int, id uuid.UUID) (int64, []repository.ListKebeleAuditLogsRow, error) {
	count, err := s.q.CountListKebeleAuditLogs(ctx, &id)
	if err != nil {
		return 0, nil, err
	}
	if count == 0 {
		return 0, nil, nil
	}
	logs, err := s.q.ListKebeleAuditLogs(ctx, repository.ListKebeleAuditLogsParams{
		ID:     &id,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	return count, logs, err
}
