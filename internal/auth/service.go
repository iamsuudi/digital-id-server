package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/iamsuudi/digital-id-server/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	_ *pgxpool.Pool
	q *sqlc.Queries
}

func NewService(db *pgxpool.Pool, q *sqlc.Queries) *Service {
	return &Service{q: q}
}

// Authenticate verifies user credentials
func (s *Service) Authenticate(ctx context.Context, email, password string) (*sqlc.User, error) {
	user, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
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

	_, err = s.q.CreateUser(ctx, sqlc.CreateUserParams{
		FirstName:  input.FirstName,
		SecondName: input.SecondName,
		LastName:   input.LastName,
		Email:      input.Email,
		Phone:      input.Phone,
		Password:   string(hashedPassword),
		Role:       input.Role, // e.g. "encoder"
	})
	return err
}

func (s *Service) StoreRefreshToken(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	_, err := s.q.CreateRefreshToken(ctx, sqlc.CreateRefreshTokenParams{
		UserID:    pgtype.Int4{Int32: int32(userID), Valid: true},
		Token:     token,
		ExpiresAt: pgtype.Timestamp{Time: expiresAt, Valid: true},
	})
	return err
}

func (s *Service) RefreshAccessToken(ctx context.Context, token string) (string, error) {
	rt, err := s.q.GetRefreshToken(ctx, token)
	if err != nil || time.Now().After(rt.ExpiresAt.Time) {
		return "", errors.New("invalid or expired refresh token")
	}

	user, err := s.q.GetUserByID(ctx, rt.UserID.Int32)
	if err != nil {
		return "", err
	}

	newJWT, err := GenerateJWT(int64(user.ID), user.Role)
	if err != nil {
		return "", err
	}

	// OPTIONAL: rotate token (invalidate old one)
	_ = s.q.DeleteRefreshToken(ctx, token)

	newRefreshToken := GenerateRandomToken(64)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	_, err = s.q.CreateRefreshToken(ctx, sqlc.CreateRefreshTokenParams{
		UserID:    pgtype.Int4{Int32: user.ID, Valid: true},
		Token:     newRefreshToken,
		ExpiresAt: pgtype.Timestamp{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return "", err
	}

	return newJWT, nil
}
