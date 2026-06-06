package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"taskgo/internal/models"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Test suite for user repository
type UserRepoTestSuite struct {
	suite.Suite
	db     *sql.DB
	repo   *PostgresUserRepository
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *UserRepoTestSuite) SetupSuite() {
	// Connect to test database
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=taskgo_db sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		s.T().Fatalf("Failed to connect to test db: %v", err)
	}

	s.db = db
	s.repo = NewPostgresUserRepository(db)

	// Create test tables
	_, err = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY,
            email VARCHAR(255) UNIQUE NOT NULL,
            password_hash VARCHAR(255) NOT NULL,
            full_name VARCHAR(255) NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		s.T().Fatalf("Failed to create test table: %v", err)
	}
}

func (s *UserRepoTestSuite) SetupTest() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 5*time.Second)

	// Clean up before each test
	_, err := s.db.ExecContext(s.ctx, "TRUNCATE users")
	if err != nil {
		s.T().Fatalf("Failed to truncate users: %v", err)
	}
}

func (s *UserRepoTestSuite) TearDownTest() {
	s.cancel()
}

func (s *UserRepoTestSuite) TearDownSuite() {
	_, err := s.db.Exec("DROP TABLE IF EXISTS users")
	if err != nil {
		s.T().Logf("Failed to drop test table: %v", err)
	}
	s.db.Close()
}

func (s *UserRepoTestSuite) TestCreateUser() {
	user := &models.User{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword123",
		FullName:     "Test User",
	}

	err := s.repo.Create(s.ctx, user)
	require.NoError(s.T(), err)
	assert.NotEmpty(s.T(), user.ID)
	assert.NotZero(s.T(), user.CreatedAt)
}

func (s *UserRepoTestSuite) TestGetUserByEmail() {
	// Create user first
	user := &models.User{
		Email:        "findme@example.com",
		PasswordHash: "hashed123",
		FullName:     "Find Me",
	}
	err := s.repo.Create(s.ctx, user)
	require.NoError(s.T(), err)

	// Find by email
	found, err := s.repo.GetByEmail(s.ctx, "findme@example.com")
	require.NoError(s.T(), err)
	assert.Equal(s.T(), user.ID, found.ID)
	assert.Equal(s.T(), user.Email, found.Email)
	assert.Equal(s.T(), user.FullName, found.FullName)
}

func (s *UserRepoTestSuite) TestGetUserByEmailNotFound() {
	_, err := s.repo.GetByEmail(s.ctx, "nonexistent@example.com")
	assert.ErrorIs(s.T(), err, ErrUserNotFound)
}

func (s *UserRepoTestSuite) TestUpdateUser() {
	// Create user
	user := &models.User{
		Email:        "update@example.com",
		PasswordHash: "oldhash",
		FullName:     "Old Name",
	}
	err := s.repo.Create(s.ctx, user)
	require.NoError(s.T(), err)

	// Update user
	user.Email = "updated@example.com"
	user.FullName = "New Name"

	err = s.repo.Update(s.ctx, user)
	require.NoError(s.T(), err)

	// Verify update
	updated, err := s.repo.GetByID(s.ctx, user.ID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "updated@example.com", updated.Email)
	assert.Equal(s.T(), "New Name", updated.FullName)
}

func (s *UserRepoTestSuite) TestDeleteUser() {
	// Create user
	user := &models.User{
		Email:        "delete@example.com",
		PasswordHash: "hash",
		FullName:     "To Delete",
	}
	err := s.repo.Create(s.ctx, user)
	require.NoError(s.T(), err)

	// Delete user
	err = s.repo.Delete(s.ctx, user.ID)
	require.NoError(s.T(), err)

	// Verify deleted
	_, err = s.repo.GetByID(s.ctx, user.ID)
	assert.ErrorIs(s.T(), err, ErrUserNotFound)
}

func TestUserRepoSuite(t *testing.T) {
	suite.Run(t, new(UserRepoTestSuite))
}
