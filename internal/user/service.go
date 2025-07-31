package user

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

func (s *Service) GetUserById(ctx context.Context, id uuid.UUID) (repository.GetUserByIDRow, error) {
	return s.q.GetUserByID(ctx, id)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (repository.GetUserByEmailRow, error) {
	return s.q.GetUserByEmail(ctx, email)
}

func (s *Service) ListUsersUnderScope(ctx context.Context, limit, offset int, query string, rank *int32, c_id, sc_id, k_id *uuid.UUID) (int64, []repository.ListUsersUnderScopeRow, error) {
	count, err := s.q.CountListUsersUnderScope(ctx, repository.CountListUsersUnderScopeParams{
		CityID:    c_id,
		SubcityID: sc_id,
		KebeleID:  k_id,
	})
	users, err := s.q.ListUsersUnderScope(ctx, repository.ListUsersUnderScopeParams{
		CityID:    c_id,
		SubcityID: sc_id,
		KebeleID:  k_id,
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	return count, users, err
}

func (s *Service) ListUsersByRole(ctx context.Context, limit, offset int, query, role_slug string) (int64, []repository.ListUsersByRoleRow, error) {
	count, err := s.q.CountListUsersByRole(ctx, role_slug)
	if err != nil {
		return 0, nil, err
	}
	if count == 0 {
		return 0, nil, nil
	}
	users, err := s.q.ListUsersByRole(ctx, repository.ListUsersByRoleParams{
		Limit:    int32(limit),
		Offset:   int32(offset),
		RoleSlug: role_slug,
	})
	if err != nil {
		return 0, nil, err
	}
	return count, users, nil
}

func (s *Service) SearchUsersUnderScope(ctx context.Context, limit, offset int, query string, rank *int32, c_id, sc_id, k_id *uuid.UUID) (int64, []repository.SearchUsersUnderScopeRow, error) {
	count, err := s.q.CountSearchUsersUnderScope(ctx, repository.CountSearchUsersUnderScopeParams{
		Query: query,
		CityID:    c_id,
		SubcityID: sc_id,
		KebeleID:  k_id,
	})
	if err != nil {
		return 0, nil, err
	}
	if count == 0 {
		return 0, nil, nil
	}

	users, err := s.q.SearchUsersUnderScope(ctx, repository.SearchUsersUnderScopeParams{
		Limit:  int32(limit),
		Offset: int32(offset),
		Query:  query,
		CityID:    c_id,
		SubcityID: sc_id,
		KebeleID:  k_id,
	})
	if err != nil {
		return count, nil, err
	}

	return count, users, nil
}

func (s *Service) SearchUsersByRole(ctx context.Context, limit, offset int, query, role_slug string) (int64, []repository.SearchUsersByRoleRow, error) {
	count, err := s.q.CountSearchUsersByRole(ctx, repository.CountSearchUsersByRoleParams{
		RoleSlug: role_slug,
		Query:    query,
	})
	if err != nil {
		return 0, nil, err
	}

	users, err := s.q.SearchUsersByRole(ctx, repository.SearchUsersByRoleParams{
		Limit:    int32(limit),
		Offset:   int32(offset),
		Query:    query,
		RoleSlug: role_slug,
	})
	if err != nil {
		return 0, nil, err
	}

	return count, users, nil
}

func (s *Service) GetEffectivePermissions(ctx context.Context, id uuid.UUID) ([]repository.GetEffectivePermissionsForUserRow, error) {
	return s.q.GetEffectivePermissionsForUser(ctx, id)
}

func (s *Service) UpdateUserRole(ctx context.Context, id uuid.UUID, role string) error {
	return s.q.UpdateUserRole(ctx, repository.UpdateUserRoleParams{
		ID: id,
		RoleSlug: role,
	})
}

func (s *Service) UpdateUserInfo(ctx context.Context, id uuid.UUID, first, second, last, email, phone string) error {
	return s.q.UpdateUserInfo(ctx, repository.UpdateUserInfoParams{
		ID: id,
		FirstName: first,
		SecondName: second,
		LastName: last,
		Email: email,
		Phone: phone,
	})
}
