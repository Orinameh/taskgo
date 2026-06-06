package service

import (
	"context"
	"taskgo/internal/auth"
	"testing"

	"taskgo/internal/models"
	"taskgo/internal/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock User Repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]models.User, int64, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]models.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	args := m.Called(ctx, id, passwordHash)
	return args.Error(0)
}

func (m *MockUserRepository) Exists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func TestUserService_Register_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	// Setup expectations
	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, repository.ErrUserNotFound)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

	// Call service
	user, token, err := service.Register(context.Background(), "test@example.com", "password123", "Test User")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.FullName)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Register_UserAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	existingUser := &models.User{
		ID:       uuid.New(),
		Email:    "existing@example.com",
		FullName: "Existing User",
	}

	mockRepo.On("GetByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)

	_, _, err := service.Register(context.Background(), "existing@example.com", "password", "New User")

	assert.Error(t, err)
	assert.Equal(t, "user already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	// Create user with known password
	password := "correctpassword"
	hashedPassword, _ := auth.HashPassword(password)

	user := &models.User{
		ID:           uuid.New(),
		Email:        "login@example.com",
		PasswordHash: hashedPassword,
		FullName:     "Login User",
	}

	mockRepo.On("GetByEmail", mock.Anything, "login@example.com").Return(user, nil)

	// Login
	returnedUser, token, err := service.Login(context.Background(), "login@example.com", password)

	assert.NoError(t, err)
	assert.NotNil(t, returnedUser)
	assert.NotEmpty(t, token)
	assert.Equal(t, user.Email, returnedUser.Email)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	hashedPassword, _ := auth.HashPassword("correctpassword")

	user := &models.User{
		ID:           uuid.New(),
		Email:        "login@example.com",
		PasswordHash: hashedPassword,
		FullName:     "Login User",
	}

	mockRepo.On("GetByEmail", mock.Anything, "login@example.com").Return(user, nil)

	_, _, err := service.Login(context.Background(), "login@example.com", "wrongpassword")

	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	mockRepo.On("GetByEmail", mock.Anything, "notfound@example.com").Return(nil, repository.ErrUserNotFound)

	_, _, err := service.Login(context.Background(), "notfound@example.com", "password")

	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
	mockRepo.AssertExpectations(t)
}
