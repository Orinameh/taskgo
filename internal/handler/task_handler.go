package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"taskgo/internal/middleware"
	"taskgo/internal/models"
	"taskgo/internal/service"
	"time"

	"github.com/google/uuid"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		parsed, err := time.Parse(time.RFC3339, *req.DueDate)
		if err == nil {
			dueDate = &parsed
		}
	}

	// claims.UserID is already uuid.UUID, use it directly
	task, err := h.taskService.CreateTask(
		r.Context(),
		claims.UserID, // Now works directly - no conversion needed
		req.Title,
		req.Description,
		req.Priority,
		dueDate,
	)

	if err != nil {
		sendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, http.StatusCreated, task)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract ID from URL path /tasks/{id}
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		sendJSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	taskID, err := uuid.Parse(pathParts[2])
	if err != nil {
		sendJSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := h.taskService.GetTask(r.Context(), taskID, claims.UserID)
	if err != nil {
		sendJSONError(w, "Task not found", http.StatusNotFound)
		return
	}

	sendJSONResponse(w, http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		sendJSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	taskID, err := uuid.Parse(pathParts[2])
	if err != nil {
		sendJSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		parsed, err := time.Parse(time.RFC3339, *req.DueDate)
		if err == nil {
			dueDate = &parsed
		}
	}

	task, err := h.taskService.UpdateTask(
		r.Context(),
		taskID,
		claims.UserID,
		req.Title,
		req.Description,
		req.Status,
		req.Priority,
		dueDate,
		req.Version,
	)

	if err != nil {
		if err.Error() == "task not found" {
			sendJSONError(w, "Task not found", http.StatusNotFound)
			return
		}
		sendJSONError(w, err.Error(), http.StatusConflict)
		return
	}

	sendJSONResponse(w, http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		sendJSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	taskID, err := uuid.Parse(pathParts[2])
	if err != nil {
		sendJSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = h.taskService.DeleteTask(r.Context(), taskID, claims.UserID)
	if err != nil {
		sendJSONError(w, "Task not found", http.StatusNotFound)
		return
	}

	sendJSONResponse(w, http.StatusOK, map[string]string{"message": "Task deleted successfully"})
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	status := r.URL.Query().Get("status")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	tasks, total, err := h.taskService.ListTasks(
		r.Context(),
		claims.UserID,
		status,
		page,
		limit,
	)

	if err != nil {
		sendJSONError(w, "Failed to list tasks", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, http.StatusOK, map[string]interface{}{
		"tasks": tasks,
		"meta": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *TaskHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		sendJSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	taskID, err := uuid.Parse(pathParts[2])
	if err != nil {
		sendJSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.taskService.AddComment(r.Context(), taskID, claims.UserID, req.Content)
	if err != nil {
		sendJSONError(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, http.StatusCreated, map[string]string{"message": "Comment added successfully"})
}

func (h *TaskHandler) GetTaskComments(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		sendJSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	taskID, err := uuid.Parse(pathParts[2])
	if err != nil {
		sendJSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Verify task belongs to user
	_, err = h.taskService.GetTask(r.Context(), taskID, claims.UserID)
	if err != nil {
		sendJSONError(w, "Task not found", http.StatusNotFound)
		return
	}

	comments, err := h.taskService.GetTaskComments(r.Context(), taskID)
	if err != nil {
		sendJSONError(w, "Failed to get comments", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, http.StatusOK, comments)
}

func sendJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.APIResponse{
		Success: true,
		Data:    data,
	})
}

func sendJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.APIResponse{
		Success: false,
		Error:   message,
	})
}
