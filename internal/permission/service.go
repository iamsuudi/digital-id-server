package permission

import (
	"context"

	"github.com/google/uuid"
	"github.com/iamsuudi/digital-id-server/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (s *Service) GetAssignablePermissions(ctx context.Context, userID uuid.UUID) ([]repository.Permission, error) {
	return s.q.GetAssignablePermissionsForActor(ctx, userID)
}
