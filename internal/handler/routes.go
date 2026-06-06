package handler

import "net/http"

// RegisterRoutes registers all API routes
func RegisterRoutes(
	userHandler *UserHandler,
	taskHandler *TaskHandler,
) {
	// User routes
	http.HandleFunc("POST /api/register", userHandler.Register)
	http.HandleFunc("POST /api/login", userHandler.Login)

	// Task routes
	http.HandleFunc("POST /api/tasks", taskHandler.CreateTask)
	http.HandleFunc("GET /api/tasks/{id}", taskHandler.GetTask)
	http.HandleFunc("PUT /api/tasks/{id}", taskHandler.UpdateTask)
	http.HandleFunc("DELETE /api/tasks/{id}", taskHandler.DeleteTask)
	http.HandleFunc("GET /api/tasks", taskHandler.ListTasks)
	http.HandleFunc("POST /api/tasks/{id}/comments", taskHandler.AddComment)
	http.HandleFunc("GET /api/tasks/{id}/comments", taskHandler.GetTaskComments)
}
