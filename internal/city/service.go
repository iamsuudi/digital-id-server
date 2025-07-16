package city

import (
	"context"

	"github.com/google/uuid"
	"github.com/iamsuudi/digital-id-server/internal/repository"
	"github.com/iamsuudi/digital-id-server/shared/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	db *pgxpool.Pool
	q  *repository.Queries
}

func NewService(dbConn *pgxpool.Pool, dbQueries *repository.Queries) *Service {
	return &Service{db: dbConn, q: dbQueries}
}

func (s *Service) CreateCity(ctx context.Context, input types.CityInput) (repository.City, error) {
	return s.q.CreateCity(ctx, input.Name)
}

func (s *Service) GetCity(ctx context.Context, id uuid.UUID) (repository.GetResidentFullRow, error) {
	return s.q.GetResidentFull(ctx, id)
}

func (s *Service) GetAll(ctx context.Context) ([]repository.ListCitiesRow, error) {
	return s.q.ListCities(ctx)
}
