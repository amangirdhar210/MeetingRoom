package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/amangirdhar210/meeting-room/internal/domain"
	"github.com/amangirdhar210/meeting-room/internal/http/dto"
	"github.com/amangirdhar210/meeting-room/internal/http/middleware"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	UserService domain.UserService
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	_, role, ok := middleware.GetUserIDRole(r.Context())
	if !ok || role != "admin" {
		http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
		return
	}

	var req dto.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	if err := h.UserService.Register(user); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.GenericResponse{Message: "user registered successfully"})
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	_, role, ok := middleware.GetUserIDRole(r.Context())
	if !ok || role != "admin" {
		http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
		return
	}

	users, err := h.UserService.GetAllUsers()
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	var resp []dto.UserDTO
	for _, u := range users {
		resp = append(resp, dto.UserDTO{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
			Role:  u.Role,
		})
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, role, ok := middleware.GetUserIDRole(r.Context())
	if !ok || role != "admin" {
		http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, `{"error":"invalid user id"}`, http.StatusBadRequest)
		return
	}
	if id == userId {
		http.Error(w, `{"error":"cannot delete yourself"}`, http.StatusForbidden)
		return
	}
	user, err := h.UserService.GetUserByID(id)
	if err != nil {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}
	if user.Email == "admin@wg.com" {
		http.Error(w, `{"error":"cannot delete superadmin"}`, http.StatusForbidden)
		return
	}

	if err := h.UserService.DeleteUserByID(id); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(dto.GenericResponse{Message: "user deleted successfully"})
}
