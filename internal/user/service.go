package user

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

func (s *Service) GetUserById(ctx context.Context, id uuid.UUID) (repository.GetUserByIDRow, error) {
	return s.q.GetUserByID(ctx, id)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (repository.GetUserByEmailRow, error) {
	return s.q.GetUserByEmail(ctx, email)
}

func (s *Service) GetAll(ctx context.Context, limit, offset int, query string) (int64, []repository.ListUsersRow, error) {
	count, _ := s.q.CountListUsers(ctx)
	users, err := s.q.ListUsers(ctx, repository.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	return count, users, err
}

func (s *Service) SearchUsers(ctx context.Context, limit, offset int, query string) (int64, []repository.SearchUsersRow, error) {
	count, _ := s.q.CountUsersSearch(ctx, query)
	users, err := s.q.SearchUsers(ctx, repository.SearchUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
		Query:  &query,
	})
	return count, users, err
}
