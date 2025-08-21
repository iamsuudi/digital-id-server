package analytics

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

func (s *Service) GetAgeGroupDistribution(ctx context.Context) ([]repository.GetAgeGroupDistributionRow, error) {
	return s.q.GetAgeGroupDistribution(ctx)
}

func (s *Service) GetGenderDistribution(ctx context.Context) ([]repository.GetGenderDistributionRow, error) {
	return s.q.GetGenderDistribution(ctx)
}

func (s *Service) GetGenderAgeGroupDistribution(ctx context.Context) ([]repository.GetGenderAgeGroupDistributionRow, error) {
	return s.q.GetGenderAgeGroupDistribution(ctx)
}
