package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"taskgo/internal/database"
	"taskgo/internal/handler"
	"taskgo/internal/repository"
	"taskgo/internal/service"
	"time"
)

func main() {
	// Initialize database
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	// NOTE: Because we removed log.Fatal from the server execution loop below,
	// this deferred db.Close() is now GUARANTEED to run cleanly on shutdown!
	defer func() {
		log.Println("Closing database connection pool...")
		db.Close()
	}()

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

	// 1. Create a dedicated, isolated routing multiplexer
	mux := http.NewServeMux()

	// Register all routes
	handler.RegisterRoutes(mux, userHandler, taskHandler)

	// Set up OS signal listening context (Listens for Ctrl+C or kill signals)
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)

	defer stop()

	// Configure the HTTP Server settings
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux, // Uses your default ServeMux registered via RegisterRoutes
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Server starting on :%s", port)

	// Spin off the server loop into a background thread (Goroutine)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server crashed unexpectedly: %v", err)
		}
	}()

	// Pause execution here until you press Ctrl+C or a termination signal is sent
	<-ctx.Done()
	log.Println("Shutdown signal received! Initiating graceful stop...")

	// Provide a strict 15-second grace period for ongoing HTTP requests to finish
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Stop accepting new traffic and wait out the active connection queue
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Graceful shutdown failed to complete: %v", err)
	}

	log.Println("Server engine stopped successfully. Cleaning up resources...")
	// The function ends here, naturally dropping into the deferred db.Close() step above!
}
