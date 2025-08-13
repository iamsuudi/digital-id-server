package audit

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

func (s *Service) GetUserById(ctx context.Context, id uuid.UUID) (repository.GetUserByIDRow, error) {
	return s.q.GetUserByID(ctx, id)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (repository.GetUserByEmailRow, error) {
	return s.q.GetUserByEmail(ctx, email)
}

func (s *Service) ListAuditLogs(ctx context.Context, limit, offset int, city, subcity, kebele *uuid.UUID) (int64, []repository.ListAuditLogsRow, error) {
	count, err := s.q.CountListAuditLogs(ctx, repository.CountListAuditLogsParams{
		CityID:    city,
		SubcityID: subcity,
		KebeleID:  kebele,
	})
	if err != nil {
		return 0, nil, err
	}
	if count == 0 {
		return 0, nil, nil
	}
	users, err := s.q.ListAuditLogs(ctx, repository.ListAuditLogsParams{
		CityID:    city,
		SubcityID: subcity,
		KebeleID:  kebele,
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	return count, users, err
}
