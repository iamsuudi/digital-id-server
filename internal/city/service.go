package city

import (
	"context"
	"fmt"

	"digital-id-server/internal/repository"
	"digital-id-server/shared/types"

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

func (s *Service) CreateCity(ctx context.Context, actorID uuid.UUID, input types.CityInput) (repository.City, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repository.City{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	city, err := qtx.CreateCity(ctx, repository.CreateCityParams{
		Name: input.Name,
		Lat:  &input.Lat,
		Lon:  &input.Lon,
	})
	if err != nil {
		return repository.City{}, err
	}

	err = qtx.InsertAuditLog(ctx, repository.InsertAuditLogParams{
		ActorID:      actorID,
		ActionType:   "CREATE_CITY",
		TargetCityID: &(city.ID),
		ObjectType:   "city",
		Diff: map[string]any{
			"after": city,
		},
	})
	if err != nil {
		return repository.City{}, err
	}

	return city, tx.Commit(ctx)
}

func (s *Service) DeleteCity(ctx context.Context, actorID, id uuid.UUID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	// 1. Before
	city, err := qtx.GetCity(ctx, id)
	if err != nil {
		return err
	}

	// 2. Delete
	err = s.q.SoftDeleteCity(ctx, id)
	if err != nil {
		return err
	}

	// 3. Insert audit log
	qtx.InsertAuditLog(ctx, repository.InsertAuditLogParams{
		ActorID:    actorID,
		ActionType: "DELETE_CITY",
		ObjectType: "city",
		Diff: map[string]any{
			"before": city,
		},
	})
	return tx.Commit(ctx)
}

func (s *Service) RemoveStaff(ctx context.Context, actorID, staffID uuid.UUID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := s.q.WithTx(tx)

	//  Before
	user, err := qtx.GetUserByID(ctx, staffID)
	if err != nil {
		return err
	}
	city, err := qtx.GetCity(ctx, *user.CityID)
	if err != nil {
		return err
	}

	// After
	err = qtx.RevokeUserPlacement(ctx, repository.RevokeUserPlacementParams{
		CityID: nil,
		ID:     staffID,
	})

	// Insert audit log
	err = qtx.InsertAuditLog(ctx, repository.InsertAuditLogParams{
		ActorID:      actorID,
		ActionType:   "REMOVE_ADMIN",
		TargetCityID: &(city.ID),
		ObjectType:   "city",
		Diff: map[string]any{
			"before": map[string]any{
				"user_name": city.AdminName,
				"city":      city.Name,
				"position":  "Admin",
			},
		},
	})

	return tx.Commit(ctx)
}

func (s *Service) AssignAdmin(ctx context.Context, actorID, cityID, adminID uuid.UUID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	city, err := qtx.GetCity(ctx, cityID)
	if err != nil {
		return err
	}

	// Remove previous admin
	if city.AdminID != nil {
		err = qtx.RevokeUserPlacement(ctx, repository.RevokeUserPlacementParams{
			CityID: nil,
			ID:     *city.AdminID,
		})
		if err != nil {
			return err
		}
	}

	err = qtx.GrantUserPlacement(ctx, repository.GrantUserPlacementParams{
		ID:     adminID,
		CityID: &cityID,
	})
	if err != nil {
		return err
	}

	user, err := qtx.GetUserByID(ctx, adminID)
	if err != nil {
		return err
	}

	// 3. Insert audit log
	err = qtx.InsertAuditLog(ctx, repository.InsertAuditLogParams{
		ActorID:      actorID,
		ActionType:   "ASSIGN_ADMIN",
		TargetCityID: &(cityID),
		ObjectType:   "city",
		Diff: map[string]any{
			"before": map[string]any{
				"user_name": city.AdminName,
				"city":      city.Name,
				"position":  "Admin",
			},
			"after": map[string]any{
				"user_name": user.FullName,
				"city":      city.Name,
				"position":  "Admin",
			},
		},
	})

	return tx.Commit(ctx)
}

func (s *Service) UpdateCityInfo(ctx context.Context, actorID, id uuid.UUID, name string, lat, lon *float64) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := s.q.WithTx(tx)

	// 1. Before
	before, err := qtx.GetCity(ctx, id)
	if err != nil {
		return err
	}

	// 2. After
	updated, err := qtx.UpdateCity(ctx, repository.UpdateCityParams{
		ID:   id,
		Name: name,
		Lat:  lat,
		Lon:  lon,
	})

	// 3. Insert audit log
	err = qtx.InsertAuditLog(ctx, repository.InsertAuditLogParams{
		ActorID:      actorID,
		TargetCityID: &id,
		ActionType:   "UPDATE_CITY",
		ObjectType:   "city",
		Diff: map[string]any{
			"before": before,
			"after":  updated,
		},
	})
	return tx.Commit(ctx)
}

func (s *Service) GetSubCities(ctx context.Context, id uuid.UUID) ([]repository.GetSubCitiesForCityRow, error) {
	return s.q.GetSubCitiesForCity(ctx, id)
}

func (s *Service) GetCity(ctx context.Context, id uuid.UUID) (repository.GetCityRow, error) {
	return s.q.GetCity(ctx, id)
}

func (s *Service) GetCities(ctx context.Context, limit, offset int) (int64, []repository.ListCitiesRow, error) {
	count, err := s.q.CountListCities(ctx)
	if err != nil {
		fmt.Println(err)
		return 0, nil, err
	}

	cities, err := s.q.ListCities(ctx, repository.ListCitiesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		fmt.Println(err)
		return 0, nil, err
	}

	return count, cities, nil
}

func (s *Service) SearchCities(ctx context.Context, limit, offset int, query string) (int64, []repository.SearchCitiesRow, error) {
	count, err := s.q.CountSearchCities(ctx, query)
	if err != nil {
		return 0, nil, err
	}

	cities, err := s.q.SearchCities(ctx, repository.SearchCitiesParams{
		Query:  query,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return 0, nil, err
	}

	return count, cities, nil
}

func (s *Service) ListAuditLogs(ctx context.Context, limit, offset int, id uuid.UUID) (int64, []repository.ListCityAuditLogsRow, error) {
	count, err := s.q.CountListCityAuditLogs(ctx, &id)
	if err != nil {
		return 0, nil, err
	}
	if count == 0 {
		return 0, nil, nil
	}
	logs, err := s.q.ListCityAuditLogs(ctx, repository.ListCityAuditLogsParams{
		ID:     &id,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	return count, logs, err
}
