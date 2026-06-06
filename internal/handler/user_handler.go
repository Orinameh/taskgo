package handler

import (
	"encoding/json"
	"net/http"
	"taskgo/internal/models"
	"taskgo/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, token, err := h.userService.Register(r.Context(), req.Email, req.Password, req.FullName)
	if err != nil {
		sendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, token, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		sendJSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	sendJSONResponse(w, http.StatusOK, map[string]interface{}{
		"user":  user,
		"token": token,
	})
}
