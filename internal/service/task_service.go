// internal/service/task_service.go
package service

import (
	"context"
	"database/sql"
	"errors"
	"taskgo/internal/models"
	"taskgo/internal/repository"
	"time"

	"github.com/google/uuid"
)

type TaskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, userID uuid.UUID, title, description string, priority int, dueDate *time.Time) (*models.Task, error) {
	task := &models.Task{
		Title:       title,
		Description: description,
		Priority:    priority,
		UserID:      userID,
	}

	if dueDate != nil {
		task.DueDate = sql.NullTime{Time: *dueDate, Valid: true}
	}

	if task.Priority == 0 {
		task.Priority = 1
	}

	err := s.taskRepo.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, taskID, userID uuid.UUID) (*models.TaskWithComments, error) {
	return s.taskRepo.GetByID(ctx, taskID, userID)
}

// New method: GetTaskComments only (without task details)
func (s *TaskService) GetTaskComments(ctx context.Context, taskID uuid.UUID) ([]models.TaskComment, error) {
	return s.taskRepo.GetComments(ctx, taskID)
}

func (s *TaskService) UpdateTask(ctx context.Context, taskID, userID uuid.UUID, title, description, status *string, priority *int, dueDate *time.Time, version int) (*models.Task, error) {
	task := &models.Task{
		ID:      taskID,
		UserID:  userID,
		Version: version,
	}

	if title != nil {
		task.Title = *title
	}
	if description != nil {
		task.Description = *description
	}
	if status != nil {
		task.Status = *status
	}
	if priority != nil {
		task.Priority = *priority
	}
	if dueDate != nil {
		task.DueDate = sql.NullTime{Time: *dueDate, Valid: true}
	}

	err := s.taskRepo.Update(ctx, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, taskID, userID uuid.UUID) error {
	return s.taskRepo.Delete(ctx, taskID, userID)
}

func (s *TaskService) ListTasks(ctx context.Context, userID uuid.UUID, status string, page, pageSize int) ([]models.Task, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.taskRepo.List(ctx, userID, status, pageSize, offset)
}

func (s *TaskService) AddComment(ctx context.Context, taskID, userID uuid.UUID, content string) error {
	if content == "" {
		return errors.New("comment cannot be empty")
	}
	return s.taskRepo.AddComment(ctx, taskID, userID, content)
}
