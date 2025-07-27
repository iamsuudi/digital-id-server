package user

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

func (s *Service) GetAllUnderScope(ctx context.Context, limit, offset int, query string) (int64, []repository.ListUsersUnderScopeRow, error) {
	count, _ := s.q.CountListUsersUnderScope(ctx, repository.CountListUsersUnderScopeParams{})
	users, err := s.q.ListUsersUnderScope(ctx, repository.ListUsersUnderScopeParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	return count, users, err
}

func (s *Service) GetAllForSuperadmin(ctx context.Context, limit, offset int, query string) (int64, []repository.ListAllUsersRow, error) {
	count, _ := s.q.CountListUsers(ctx)
	users, err := s.q.ListAllUsers(ctx, repository.ListAllUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	return count, users, err
}

func (s *Service) GetByRole(ctx context.Context, limit, offset int, query, role_slug string) (int64, []repository.ListByRoleRow, error) {
	count, _ := s.q.CountListByRole(ctx, role_slug)
	users, err := s.q.ListByRole(ctx, repository.ListByRoleParams{
		Limit:  int32(limit),
		Offset: int32(offset),
		RoleSlug: role_slug,
	})
	return count, users, err
}

func (s *Service) SearchUsersUnderScope(ctx context.Context, limit, offset int, query string) (int64, []repository.SearchUsersUnderScopeRow, error) {
	count, err := s.q.CountUsersSearchUnderScope(ctx, repository.CountUsersSearchUnderScopeParams{
		Query: query,
	})
	if err != nil {
		return 0, nil, err
	}

	users, err := s.q.SearchUsersUnderScope(ctx, repository.SearchUsersUnderScopeParams{
		Limit:  int32(limit),
		Offset: int32(offset),
		Query:  query,
	})
	if err != nil {
		return 0, nil, err
	}

	return count, users, nil
}

func (s *Service) SearchUsersForSuperadmin(ctx context.Context, limit, offset int, query string) (int64, []repository.SearchUsersRow, error) {
	count, err := s.q.CountUsersSearch(ctx, query)
	if err != nil {
		return 0, nil, err
	}

	users, err := s.q.SearchUsers(ctx, repository.SearchUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
		Query:  query,
	})
	if err != nil {
		return 0, nil, err
	}

	return count, users, nil
}

func (s *Service) SearchByRole(ctx context.Context, limit, offset int, query, role_slug string) (int64, []repository.SearchByRoleRow, error) {
	count, err := s.q.CountByRoleSearch(ctx, repository.CountByRoleSearchParams{
		RoleSlug: role_slug,
		Query: query,
	})
	if err != nil {
		return 0, nil, err
	}

	users, err := s.q.SearchByRole(ctx, repository.SearchByRoleParams{
		Limit:  int32(limit),
		Offset: int32(offset),
		Query:  query,
		RoleSlug: role_slug,
	})
	if err != nil {
		return 0, nil, err
	}

	return count, users, nil
}
