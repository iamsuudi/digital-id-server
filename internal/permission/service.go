package permission

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

func (s *Service) GetAllPermissions(ctx context.Context) ([]repository.Permission, error) {
	return s.q.ListPermissions(ctx)
}

func (s *Service) GetAssignablePermissionsForActor(ctx context.Context, id uuid.UUID) ([]repository.Permission, error) {
	return s.q.GetAssignablePermissionsForActor(ctx, id)
}

func (s *Service) GetUniversalPermissionsForUser(ctx context.Context, actor, target string, id uuid.UUID) ([]repository.GetUniversalPermissionMatrixForUserRow, error) {
	return s.q.GetUniversalPermissionMatrixForUser(ctx, repository.GetUniversalPermissionMatrixForUserParams{
		ActorRoleSlug: actor,
		TargetRoleSlug: target,
		TargetUserID: id,
	})
}

func (s *Service) OverrideUserPermission(ctx context.Context, actor, target uuid.UUID, permission string, granted bool) error {
	return s.q.SetUserPermissionOverride(ctx, repository.SetUserPermissionOverrideParams{
		UserID: target,
		GrantedBy: &actor,
		IsGranted: granted,
		PermissionName: permission,
	})
}

func (s *Service) RemoveUserPermissionOverride(ctx context.Context, target uuid.UUID, permission string) error {
	return s.q.RemoveUserPermissionOverride(ctx, repository.RemoveUserPermissionOverrideParams{
		UserID: target,
		PermissionName: permission,
	})
}
