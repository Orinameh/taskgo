package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"taskgo/internal/models"

	"github.com/google/uuid"
)

// Custom errors
var (
	ErrUserNotFound = errors.New("user not found")
	ErrTaskNotFound = errors.New("task not found")
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (id, email, password_hash, full_name)
        VALUES ($1, $2, $3, $4)
        RETURNING created_at, updated_at
    `

	user.ID = uuid.New()
	err := r.db.QueryRowContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FullName,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
        SELECT id, email, password_hash, full_name, created_at, updated_at
        FROM users
        WHERE email = $1
    `

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
        SELECT id, email, password_hash, full_name, created_at, updated_at
        FROM users
        WHERE id = $1
    `

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
        UPDATE users
        SET email = $1, full_name = $2, updated_at = CURRENT_TIMESTAMP
        WHERE id = $3
        RETURNING updated_at
    `

	err := r.db.QueryRowContext(ctx, query, user.Email, user.FullName, user.ID).
		Scan(&user.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}
