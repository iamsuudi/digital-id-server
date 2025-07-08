package user

import (
	"context"
	generated "github.com/iamsuudi/digital-id-server/db/generated"
)

type Service struct {
	q *generated.Queries
}

func NewService(q *generated.Queries) *Service {
	return &Service{q}
}

func (s *Service) CreateUser(ctx context.Context, name, email string) (*generated.User, error) {
	u, err := s.q.CreateUser(ctx, generated.CreateUserParams{
		Name:  name,
		Email: email,
	})
	if err != nil {
		return nil, err
	}
	return &u, nil
}
