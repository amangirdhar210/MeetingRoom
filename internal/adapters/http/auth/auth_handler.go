package auth

import (
	"encoding/json"
	"net/http"

	"github.com/amangirdhar210/meeting-room/internal/adapters/httputil"
	"github.com/amangirdhar210/meeting-room/internal/core/service"
	"github.com/amangirdhar210/meeting-room/internal/http/dto"
)

type Handler struct {
	authService service.AuthService
}

func NewHandler(authService service.AuthService) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		httputil.HandleError(w, err)
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

	httputil.RespondWithJSON(w, http.StatusOK, resp)
}
