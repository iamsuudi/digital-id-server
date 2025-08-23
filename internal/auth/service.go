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
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	user, err := qtx.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email")
	}

	account, err := qtx.GetAccount(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return &user, tx.Commit(ctx)
}

// RegisterUser creates a new user account with hashed password.
func (s *Service) RegisterUser(ctx context.Context, input types.UserRegisterInput) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	_, err = qtx.GetUserByEmail(ctx, input.Email)
	if err == nil {
		return errors.New("Email already exists!")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := qtx.CreateUser(ctx, repository.CreateUserParams{
		FirstName:  input.FirstName,
		SecondName: input.SecondName,
		LastName:   input.LastName,
		Email:      input.Email,
		Phone:      input.Phone,
		RoleSlug:   input.RoleSlug,
	})
	if err != nil {
		return err
	}

	err = qtx.CreateAccount(ctx, repository.CreateAccountParams{
		UserID:       user.ID,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
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
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return "", "", err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	rt, err := qtx.GetRefreshToken(ctx, token)
	if err != nil || time.Now().After(rt.ExpiresAt) {
		return "", "", errors.New("invalid or expired refresh token")
	}

	user, err := qtx.GetUserByID(ctx, rt.UserID)
	if err != nil {
		return "", "", err
	}

	newJWT, err := GenerateJWT(user.ID, user.RoleSlug)
	if err != nil {
		return "", "", err
	}

	// Rotate refresh token (optional, but more secure)
	_ = qtx.DeleteRefreshToken(ctx, token)

	newRefreshToken := GenerateRandomToken(64)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	_, err = qtx.CreateRefreshToken(ctx, repository.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return "", "", err
	}

	return newJWT, newRefreshToken, tx.Commit(ctx)
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
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	// Get user by email
	user, err := qtx.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	// Generate secure token
	token := GenerateRandomToken(32)

	// Set expiration time (1 hour from now)
	expiresAt := time.Now().Add(30 * time.Minute)

	// Create password reset token
	_, err = qtx.CreatePasswordResetToken(ctx, repository.CreatePasswordResetTokenParams{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return "", err
	}

	return token, tx.Commit(ctx)
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
		UserID:       resetToken.UserID,
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
