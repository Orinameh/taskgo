package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Custom UUID type that uses v7
type UUIDv7 uuid.UUID

// Generate new UUID v7
func NewUUIDv7() UUIDv7 {
	// github.com/google/uuid doesn't support v7 yet
	// Use a different library or implement manually
	return UUIDv7(uuid.New()) // Placeholder - will implement v7 properly
}

// Or just use standard uuid with v7 from other library
type User struct {
	ID           uuid.UUID `json:"id"` // Will generate v7
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Task struct {
	ID          uuid.UUID    `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      string       `json:"status"`
	Priority    int          `json:"priority"`
	DueDate     sql.NullTime `json:"due_date,omitempty"`
	UserID      uuid.UUID    `json:"user_id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at,omitempty"`
	Version     int          `json:"version"`
}

type TaskWithComments struct {
	Task
	Comments []TaskComment `json:"comments,omitempty"`
}

type TaskComment struct {
	ID        int64     `json:"id"`
	TaskID    int64     `json:"task_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type AuditLog struct {
	ID           int64           `json:"id"`
	UserID       sql.NullInt64   `json:"user_id,omitempty"`
	Action       string          `json:"action"`
	ResourceType string          `json:"resource_type"`
	ResourceID   sql.NullInt64   `json:"resource_id,omitempty"`
	Details      json.RawMessage `json:"details,omitempty"`
	IPAddress    string          `json:"ip_address,omitempty"`
	UserAgent    string          `json:"user_agent,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
}

type CreateTaskRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Priority    int     `json:"priority"`
	DueDate     *string `json:"due_date,omitempty"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
	Priority    *int    `json:"priority,omitempty"`
	DueDate     *string `json:"due_date,omitempty"`
	Version     int     `json:"version"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}
