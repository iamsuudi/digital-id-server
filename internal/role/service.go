package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"digital-id-server/internal/repository"
)

type Service struct {
	db *pgxpool.Pool
	q  *repository.Queries
}

func NewService(dbConn *pgxpool.Pool, dbQueries *repository.Queries) *Service {
	return &Service{db: dbConn, q: dbQueries}
}

func (s *Service) GetRoles(ctx context.Context) ([]repository.Role, error) {
	return s.q.ListRoles(ctx)
}

func (s *Service) GetRolesTree(ctx context.Context, role string) ([]repository.GetRoleTreeRow, error) {
	return s.q.GetRoleTree(ctx, role)
}

func (s *Service) GetRolesPermissions(ctx context.Context) ([]repository.ListRolePermissionsRow, error) {
	return s.q.ListRolePermissions(ctx)
}

func (s *Service) GetMyRoleLevelRank(ctx context.Context, id uuid.UUID) (int32, error) {
	return s.q.GetCurrentUserMaxRoleLevel(ctx, id)
}

func (s *Service) CanManipulateRole(ctx context.Context, id uuid.UUID, role_slug string) (bool, error) {
	return s.q.CanActorManipulateRole(ctx, repository.CanActorManipulateRoleParams{
		ID:   id,
		Role: role_slug,
	})
}

func (s *Service) CanGrantPermissionToRole(ctx context.Context, id uuid.UUID, role string, permission string) (bool, error) {
	return s.q.CanActorGrantPermissionToRole(ctx, repository.CanActorGrantPermissionToRoleParams{
		ID:         id,
		Role:       role,
		Permission: permission,
	})
}

func (s *Service) GrantPermissionToRole(ctx context.Context, role string, permission string) error {
	return s.q.GrantPermissionToRole(ctx, repository.GrantPermissionToRoleParams{
		RoleSlug:       role,
		PermissionName: permission,
	})
}

func (s *Service) RevokePermissionFromRole(ctx context.Context, role string, permission string) error {
	return s.q.RevokePermissionFromRole(ctx, repository.RevokePermissionFromRoleParams{
		RoleSlug:       role,
		PermissionName: permission,
	})
}
