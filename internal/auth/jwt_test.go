package auth

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password := "mysecretpassword123"

	hash, err := HashPassword(password)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)

	// Test valid password
	valid := CheckPasswordHash(password, hash)
	assert.True(t, valid)

	// Test invalid password
	invalid := CheckPasswordHash("wrongpassword", hash)
	assert.False(t, invalid)
}

func TestGenerateJWT(t *testing.T) {
	// Set JWT secret for testing
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	userID := uuid.New()
	email := "test@example.com"

	token, err := GenerateJWT(userID, email)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateJWT(t *testing.T) {
	// Set JWT secret for testing
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	userID := uuid.New()
	email := "test@example.com"

	// Generate token
	token, err := GenerateJWT(userID, email)
	require.NoError(t, err)

	// Validate token
	claims, err := ValidateJWT(token)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.NotNil(t, claims.ExpiresAt)
}

func TestValidateExpiredJWT(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	// Create expired claims
	claims := &Claims{
		UserID: uuid.New(),
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	require.NoError(t, err)

	// Validate should fail
	_, err = ValidateJWT(tokenString)
	assert.Error(t, err)
}

func TestValidateInvalidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	invalidToken := "invalid.token.string"

	_, err := ValidateJWT(invalidToken)
	assert.Error(t, err)
}
