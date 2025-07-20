package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iamsuudi/digital-id-server/internal/repository"
	"github.com/iamsuudi/digital-id-server/shared/types"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db *pgxpool.Pool
	q  repository.Querier
}

func NewService(dbConn *pgxpool.Pool, dbQueries repository.Querier) *Service {
	return &Service{db: dbConn, q: dbQueries}
}

// Authenticate verifies a user's email and password.
func (s *Service) Authenticate(ctx context.Context, email, password string) (*repository.GetUserByEmailRow, error) {
	user, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return &user, nil
}

// RegisterUser creates a new user account with hashed password.
func (s *Service) RegisterUser(ctx context.Context, input types.UserRegisterInput) error {
	_, err := s.q.GetUserByEmail(ctx, input.Email)
	if err == nil {
		return fmt.Errorf("email %s already in use", input.Email)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	_, err = s.q.CreateUser(ctx, repository.CreateUserParams{
		FirstName:    input.FirstName,
		SecondName:   input.SecondName,
		LastName:     input.LastName,
		Email:        input.Email,
		Phone:        input.Phone,
		PasswordHash: string(hashedPassword),
		RoleSlug:     input.RoleSlug,
	})
	return err
}

// StoreRefreshToken saves a refresh token in the database.
func (s *Service) StoreRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	_, err := s.q.CreateRefreshToken(ctx, repository.CreateRefreshTokenParams{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	})
	return err
}

// RefreshAccessToken validates a refresh token and returns a new JWT.
func (s *Service) RefreshAccessToken(ctx context.Context, token string) (string, error) {
	rt, err := s.q.GetRefreshToken(ctx, token)
	if err != nil || time.Now().After(rt.ExpiresAt) {
		return "", errors.New("invalid or expired refresh token")
	}

	user, err := s.q.GetUserByID(ctx, rt.UserID)
	if err != nil {
		return "", err
	}

	newJWT, err := GenerateJWT(user.ID, user.RoleSlug)
	if err != nil {
		return "", err
	}

	// Rotate refresh token (optional, but more secure)
	_ = s.q.DeleteRefreshToken(ctx, token)

	newRefreshToken := GenerateRandomToken(64)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	_, err = s.q.CreateRefreshToken(ctx, repository.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return "", err
	}

	return newJWT, nil
}

// DeleteRefreshToken removes a refresh token from the database.
func (s *Service) DeleteRefreshToken(ctx context.Context, token string) error {
	return s.q.DeleteRefreshToken(ctx, token)
}

func (s *Service) GetUserByID(ctx context.Context, userID uuid.UUID) (*repository.GetUserByIDRow, error) {
	user, err := s.q.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
