package test

import (
	"database/sql"
	"fmt"
	"os"
	"taskgo/internal/models"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// SetupTestDB creates a test database connection
func SetupTestDB(t *testing.T) *sql.DB {
	// Use test database. Always create a separate one
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env:", err)
	}

	var connStr string

	// If DATABASE_URL is provided, use it
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		connStr = databaseURL
	} else {
		// Fall back to individual environment variables
		connStr = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Clean up before test
	cleanupTestDB(db)

	return db
}

func cleanupTestDB(db *sql.DB) {
	db.Exec("TRUNCATE users, tasks, task_comments CASCADE")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GenerateTestUser creates a test user
func GenerateTestUser(db *sql.DB) (*models.User, error) {
	// Implementation
	return nil, nil
}
