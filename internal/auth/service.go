package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"digital-id-server/internal/repository"
	"digital-id-server/shared/types"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db *pgxpool.Pool
	q  *repository.Queries
}

func NewService(dbConn *pgxpool.Pool, dbQueries *repository.Queries) *Service {
	return &Service{db: dbConn, q: dbQueries}
}

// Authenticate verifies a user's email and password.
func (s *Service) Authenticate(ctx context.Context, email, password string) (*repository.GetUserByEmailRow, error) {
	user, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
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
func (s *Service) RefreshAccessToken(ctx context.Context, token string) (string, string, error) {
	rt, err := s.q.GetRefreshToken(ctx, token)
	if err != nil || time.Now().After(rt.ExpiresAt) {
		return "", "", errors.New("invalid or expired refresh token")
	}

	user, err := s.q.GetUserByID(ctx, rt.UserID)
	if err != nil {
		return "", "", err
	}

	newJWT, err := GenerateJWT(user.ID, user.RoleSlug)
	if err != nil {
		return "", "", err
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
		return "", "", err
	}

	return newJWT, newRefreshToken, nil
}

// DeleteRefreshToken removes a refresh token from the database.
func (s *Service) DeleteRefreshToken(ctx context.Context, token string) error {
	return s.q.DeleteRefreshToken(ctx, token)
}

func (s *Service) GetUserByID(ctx context.Context, userID uuid.UUID) (repository.GetUserByIDRow, error) {
	user, err := s.q.GetUserByID(ctx, userID)
	if err != nil {
		return repository.GetUserByIDRow{}, err
	}
	return user, nil
}

func (s *Service) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	// Get user by email
	user, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	// Generate secure token
	token := GenerateRandomToken(32)

	// Set expiration time (1 hour from now)
	expiresAt := time.Now().Add(30 * time.Minute)

	// Create password reset token
	_, err = s.q.CreatePasswordResetToken(ctx, repository.CreatePasswordResetTokenParams{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	})
	return token, err
}

func (s *Service) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Get valid token
	resetToken, err := s.q.GetValidPasswordResetToken(ctx, token)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	// Update user password
	err = qtx.UpdateUserPassword(ctx, repository.UpdateUserPasswordParams{
		PasswordHash: string(hashedPassword),
		ID:           resetToken.UserID,
	})
	if err != nil {
		return err
	}

	// Mark token as used
	err = qtx.MarkTokenAsUsed(ctx, resetToken.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Service) CleanupExpiredTokens(ctx context.Context) error {
	return s.q.DeleteExpiredTokens(ctx)
}
