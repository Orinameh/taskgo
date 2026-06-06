package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"taskgo/internal/models"

	"github.com/google/uuid"
)

type PostgresTaskRepository struct {
	db *sql.DB
}

func NewPostgresTaskRepository(db *sql.DB) *PostgresTaskRepository {
	return &PostgresTaskRepository{db: db}
}

func (r *PostgresTaskRepository) Create(ctx context.Context, task *models.Task) error {
	query := `
        INSERT INTO tasks (id, title, description, priority, due_date, user_id, version)
        VALUES ($1, $2, $3, $4, $5, $6, 1)
        RETURNING created_at, updated_at
    `

	task.ID = uuid.New()
	err := r.db.QueryRowContext(ctx, query,
		task.ID,
		task.Title,
		task.Description,
		task.Priority,
		task.DueDate,
		task.UserID,
	).Scan(&task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

func (r *PostgresTaskRepository) GetByID(ctx context.Context, id, userID uuid.UUID) (*models.TaskWithComments, error) {
	// Get task
	taskQuery := `
        SELECT id, title, description, status, priority, due_date, user_id, 
               created_at, updated_at, version
        FROM tasks
        WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
    `

	task := &models.Task{}
	var dueDate sql.NullTime
	err := r.db.QueryRowContext(ctx, taskQuery, id, userID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&dueDate,
		&task.UserID,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	task.DueDate = dueDate

	// Get comments
	commentsQuery := `
        SELECT id, task_id, user_id, content, created_at
        FROM task_comments
        WHERE task_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, commentsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	var comments []models.TaskComment
	for rows.Next() {
		var comment models.TaskComment
		err := rows.Scan(
			&comment.ID,
			&comment.TaskID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	return &models.TaskWithComments{
		Task:     *task,
		Comments: comments,
	}, nil
}

func (r *PostgresTaskRepository) Update(ctx context.Context, task *models.Task) error {
	query := `
        UPDATE tasks
        SET title = COALESCE(NULLIF($1, ''), title),
            description = COALESCE(NULLIF($2, ''), description),
            status = COALESCE(NULLIF($3, ''), status),
            priority = COALESCE($4, priority),
            due_date = COALESCE($5, due_date),
            updated_at = CURRENT_TIMESTAMP,
            version = version + 1
        WHERE id = $6 AND user_id = $7 AND version = $8 AND deleted_at IS NULL
        RETURNING version, updated_at
    `

	var newVersion int
	var updatedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.DueDate,
		task.ID,
		task.UserID,
		task.Version,
	).Scan(&newVersion, &updatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTaskNotFound
		}
		return fmt.Errorf("failed to update task: %w", err)
	}

	task.Version = newVersion
	if updatedAt.Valid {
		task.UpdatedAt = updatedAt.Time
	}

	return nil
}

func (r *PostgresTaskRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	query := `
        UPDATE tasks
        SET deleted_at = CURRENT_TIMESTAMP
        WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
        RETURNING id
    `

	var deletedID uuid.UUID
	err := r.db.QueryRowContext(ctx, query, id, userID).Scan(&deletedID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTaskNotFound
		}
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

func (r *PostgresTaskRepository) List(ctx context.Context, userID uuid.UUID, status string, limit, offset int) ([]models.Task, int64, error) {
	// Count total
	countQuery := `
        SELECT COUNT(*)
        FROM tasks
        WHERE user_id = $1 AND deleted_at IS NULL
    `
	args := []interface{}{userID}

	if status != "" {
		countQuery += " AND status = $2"
		args = append(args, status)
	}

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tasks: %w", err)
	}

	// Get tasks
	listQuery := `
        SELECT id, title, description, status, priority, due_date, user_id, 
               created_at, updated_at, version
        FROM tasks
        WHERE user_id = $1 AND deleted_at IS NULL
    `
	listArgs := []interface{}{userID}
	argIndex := 2

	if status != "" {
		listQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		listArgs = append(listArgs, status)
		argIndex++
	}

	listQuery += fmt.Sprintf(" ORDER BY priority DESC, due_date ASC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	listArgs = append(listArgs, limit, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var dueDate sql.NullTime

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&dueDate,
			&task.UserID,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.Version,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan task: %w", err)
		}
		task.DueDate = dueDate
		tasks = append(tasks, task)
	}

	return tasks, total, nil
}

func (r *PostgresTaskRepository) AddComment(ctx context.Context, taskID, userID uuid.UUID, content string) error {
	// Verify task exists
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM tasks WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL)`
	err := r.db.QueryRowContext(ctx, checkQuery, taskID, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to verify task: %w", err)
	}

	if !exists {
		return ErrTaskNotFound
	}

	// Add comment
	query := `
        INSERT INTO task_comments (id, task_id, user_id, content)
        VALUES ($1, $2, $3, $4)
    `

	commentID := uuid.New()
	_, err = r.db.ExecContext(ctx, query, commentID, taskID, userID, content)
	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	return nil
}

// Add this method to your PostgresTaskRepository struct

// GetComments retrieves all comments for a specific task
func (r *PostgresTaskRepository) GetComments(ctx context.Context, taskID uuid.UUID) ([]models.TaskComment, error) {
	query := `
        SELECT id, task_id, user_id, content, created_at
        FROM task_comments
        WHERE task_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	var comments []models.TaskComment
	for rows.Next() {
		var comment models.TaskComment
		err := rows.Scan(
			&comment.ID,
			&comment.TaskID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
