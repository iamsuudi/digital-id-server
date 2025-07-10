package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

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

// Authenticate verifies user credentials
func (s *Service) Authenticate(ctx context.Context, email, password string) (*repository.Users, error) {
	user, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return &user, nil
}

// RegisterUser registers a new user (encoder, superadmin, etc.)
func (s *Service) RegisterUser(ctx context.Context, input RegisterInput) error {
	// Check if user already exists
	_, err := s.q.GetUserByEmail(ctx, input.Email)
	if err == nil {
		return fmt.Errorf("email %s already in use", input.Email)
	}
	// if err != sql.ErrNoRows {
	//     return err // any other db error
	// }

	// Hash password
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
		Role:         input.Role, // e.g. "encoder"
	})
	return err
}

func (s *Service) StoreRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	_, err := s.q.CreateRefreshToken(ctx, repository.CreateRefreshTokenParams{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	})
	return err
}

func (s *Service) RefreshAccessToken(ctx context.Context, token string) (string, error) {
	rt, err := s.q.GetRefreshToken(ctx, token)
	if err != nil || time.Now().After(rt.ExpiresAt) {
		return "", errors.New("invalid or expired refresh token")
	}

	user, err := s.q.GetUserByID(ctx, rt.UserID)
	if err != nil {
		return "", err
	}

	newJWT, err := GenerateJWT(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	// OPTIONAL: rotate token (invalidate old one)
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
