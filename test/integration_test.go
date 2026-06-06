package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"taskgo/internal/database"
	"taskgo/internal/handler"
	"testing"

	"taskgo/internal/models"
	"taskgo/internal/repository"
	"taskgo/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer() *httptest.Server {
	// Setup real database for integration tests
	db, err := database.NewPostgresDB()
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}

	userRepo := repository.NewPostgresUserRepository(db.DB)
	taskRepo := repository.NewPostgresTaskRepository(db.DB)

	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo)

	userHandler := handler.NewUserHandler(userService)
	taskHandler := handler.NewTaskHandler(taskService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/register", userHandler.Register)
	mux.HandleFunc("POST /api/login", userHandler.Login)
	mux.HandleFunc("POST /api/tasks", taskHandler.CreateTask)

	return httptest.NewServer(mux)
}

func TestIntegration_RegisterAndLogin(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Register user
	registerReq := models.RegisterRequest{
		Email:    "integration@example.com",
		Password: "testpass123",
		FullName: "Integration Test",
	}

	body, _ := json.Marshal(registerReq)
	resp, err := http.Post(server.URL+"/api/register", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Parse response
	var registerResp struct {
		Data struct {
			User  models.User `json:"user"`
			Token string      `json:"token"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&registerResp)
	require.NoError(t, err)
	assert.NotEmpty(t, registerResp.Data.Token)

	// Login with same credentials
	loginReq := models.LoginRequest{
		Email:    "integration@example.com",
		Password: "testpass123",
	}

	body, _ = json.Marshal(loginReq)
	resp, err = http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
