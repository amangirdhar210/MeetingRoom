package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/amangirdhar210/meeting-room/internal/domain"
	"github.com/amangirdhar210/meeting-room/internal/http/dto"
)

type AuthHandler struct {
	AuthService domain.AuthService
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	token, user, err := h.AuthService.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusUnauthorized)
		return
	}

	resp := dto.LoginUserResponse{
		Token: token,
		User: dto.UserDTO{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}

	json.NewEncoder(w).Encode(resp)
}
