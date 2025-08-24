package permission

import (
	"context"

	"digital-id-server/internal/repository"

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

func (s *Service) CreatePermission(ctx context.Context, name, label, description string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := s.q.WithTx(tx)

	err = qtx.CreatePermission(ctx, repository.CreatePermissionParams{
		Name:        name,
		Label:       label,
		Description: &description,
	})
	if err != nil {
		return err
	}

	err = qtx.GrantPermissionToRole(ctx, repository.GrantPermissionToRoleParams{
		RoleSlug:       "superadmin",
		PermissionName: name,
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Service) DeletePermission(ctx context.Context, name string) error {
	return s.q.DeletePermission(ctx, name)
}

func (s *Service) GetAllPermissions(ctx context.Context) ([]repository.Permission, error) {
	return s.q.ListPermissions(ctx)
}

func (s *Service) GetAssignablePermissionsForActor(ctx context.Context, id uuid.UUID) ([]repository.Permission, error) {
	return s.q.GetAssignablePermissionsForActor(ctx, id)
}

func (s *Service) GetUniversalPermissionsForUser(ctx context.Context, actor, target string, id uuid.UUID) ([]repository.GetUniversalPermissionMatrixForUserRow, error) {
	return s.q.GetUniversalPermissionMatrixForUser(ctx, repository.GetUniversalPermissionMatrixForUserParams{
		ActorRoleSlug:  actor,
		TargetRoleSlug: target,
		TargetUserID:   id,
	})
}

func (s *Service) OverrideUserPermission(ctx context.Context, actor, target uuid.UUID, permission string, granted bool) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := s.q.WithTx(tx)

	// 1. Before
	before, err := qtx.GetPermissionOverridesForUser(ctx, target)
	if err != nil {
		return err
	}

	// 2. After
	err = qtx.SetUserPermissionOverride(ctx, repository.SetUserPermissionOverrideParams{
		UserID:         target,
		GrantedBy:      actor,
		IsGranted:      granted,
		PermissionName: permission,
	})
	if err != nil {
		return err
	}

	after, err := qtx.GetPermissionOverridesForUser(ctx, target)
	if err != nil {
		return err
	}

	// 3. Insert audit log
	err = qtx.InsertAuditLog(ctx, repository.InsertAuditLogParams{
		ActorID:      actor,
		ActionType:   "OVERRIDE_PERMISSION",
		TargetUserID: &target,
		ObjectType:   "user",
		Diff: map[string]any{
			"before": before,
			"after":  after,
		},
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Service) RemoveUserPermissionOverride(ctx context.Context, actor, target uuid.UUID, permission string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := s.q.WithTx(tx)

	// 1. Before
	before, err := qtx.GetPermissionOverridesForUser(ctx, target)
	if err != nil {
		return err
	}

	// 2. After
	err = qtx.RemoveUserPermissionOverride(ctx, repository.RemoveUserPermissionOverrideParams{
		UserID:         target,
		PermissionName: permission,
	})
	if err != nil {
		return err
	}

	after, err := qtx.GetPermissionOverridesForUser(ctx, target)
	if err != nil {
		return err
	}

	// 3. Insert audit log
	err = qtx.InsertAuditLog(ctx, repository.InsertAuditLogParams{
		ActorID:      actor,
		ActionType:   "RESET_OVERRIDEN_PERMISSION",
		TargetUserID: &target,
		ObjectType:   "user",
		Diff: map[string]any{
			"before": before,
			"after":  after,
		},
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
