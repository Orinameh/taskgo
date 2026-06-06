package main

import (
	"log"
	"net/http"
	"os"
	"taskgo/internal/database"
	"taskgo/internal/handler"
	"taskgo/internal/repository"
	"taskgo/internal/service"
)

func main() {
	// Initialize database
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewPostgresUserRepository(db.DB)
	taskRepo := repository.NewPostgresTaskRepository(db.DB)

	// Initialize services
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	taskHandler := handler.NewTaskHandler(taskService)

	// Setup routes
	// http.HandleFunc("POST /api/register", userHandler.Register)
	// http.HandleFunc("POST /api/login", userHandler.Login)
	// http.HandleFunc("POST /api/tasks", taskHandler.CreateTask)
	// http.HandleFunc("GET /api/tasks/{id}", taskHandler.GetTask)
	// http.HandleFunc("PUT /api/tasks/{id}", taskHandler.UpdateTask)
	// http.HandleFunc("DELETE /api/tasks/{id}", taskHandler.DeleteTask)
	// http.HandleFunc("GET /api/tasks", taskHandler.ListTasks)
	// http.HandleFunc("POST /api/tasks/{id}/comments", taskHandler.AddComment)
	// http.HandleFunc("GET /api/tasks/{id}/comments", taskHandler.GetTaskComments)

	// Register all routes
	handler.RegisterRoutes(userHandler, taskHandler)

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
