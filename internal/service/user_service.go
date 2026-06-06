package service

import (
	"context"
	"errors"
	"taskgo/internal/auth"
	"taskgo/internal/models"
	"taskgo/internal/repository"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Register(ctx context.Context, email, password, fullName string) (*models.User, string, error) {
	// Check if user exists
	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return nil, "", err
	}
	if existing != nil {
		return nil, "", errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	// Create user
	user := &models.User{
		Email:        email,
		PasswordHash: hashedPassword,
		FullName:     fullName,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, "", err
	}

	// Generate token - now works with uuid.UUID
	token, err := auth.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	// Get user
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, "", errors.New("invalid credentials")
		}
		return nil, "", err
	}

	// Check password
	if !auth.CheckPasswordHash(password, user.PasswordHash) {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate token - now works with uuid.UUID
	token, err := auth.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}
