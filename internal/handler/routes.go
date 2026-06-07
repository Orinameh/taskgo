package handler

import "net/http"

// RegisterRoutes registers all API routes
func RegisterRoutes(
	mux *http.ServeMux,
	userHandler *UserHandler,
	taskHandler *TaskHandler,
) {

	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// User routes
	mux.HandleFunc("POST /api/register", userHandler.Register)
	mux.HandleFunc("POST /api/login", userHandler.Login)

	// Task routes
	mux.HandleFunc("POST /api/tasks", taskHandler.CreateTask)
	mux.HandleFunc("GET /api/tasks/{id}", taskHandler.GetTask)
	mux.HandleFunc("PUT /api/tasks/{id}", taskHandler.UpdateTask)
	mux.HandleFunc("DELETE /api/tasks/{id}", taskHandler.DeleteTask)
	mux.HandleFunc("GET /api/tasks", taskHandler.ListTasks)
	mux.HandleFunc("POST /api/tasks/{id}/comments", taskHandler.AddComment)
	mux.HandleFunc("GET /api/tasks/{id}/comments", taskHandler.GetTaskComments)
}
