package auth_test

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/iamsuudi/digital-id-server/internal/auth"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Set the JWT_SECRET environment variable for test
	_ = os.Setenv("JWT_SECRET", "testsecret")
}

func TestGenerateAndParseJWT(t *testing.T) {
	userID := uuid.New()
	role := "encoder"

	token, err := auth.GenerateJWT(userID, role)
	assert.NoError(t, err, "should generate token without error")
	assert.NotEmpty(t, token, "token should not be empty")

	claims, err := auth.ParseJWT(token)
	assert.NoError(t, err, "should parse token without error")
	assert.Equal(t, userID, claims.UserID, "user ID should match")
	assert.Equal(t, role, claims.Role, "role should match")
	assert.WithinDuration(t, time.Now(), claims.IssuedAt.Time, 2*time.Second, "issuedAt should be recent")
}

func TestParseInvalidJWT(t *testing.T) {
	_, err := auth.ParseJWT("invalid.token.here")
	assert.Error(t, err, "should fail to parse invalid token")
}
