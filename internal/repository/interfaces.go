package repository

import (
	"context"
	"taskgo/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// TaskRepository defines all task-related database operations
type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id, userID uuid.UUID) (*models.TaskWithComments, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
	List(ctx context.Context, userID uuid.UUID, status string, limit, offset int) ([]models.Task, int64, error)
	AddComment(ctx context.Context, taskID, userID uuid.UUID, content string) error
	GetComments(ctx context.Context, taskID uuid.UUID) ([]models.TaskComment, error)
}
