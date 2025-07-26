package random

import (
	"context"

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

func (s *Service) GetAll(ctx context.Context) ([]*repository.User, error) {
	return s.q.GetAllUsers(ctx)
}
