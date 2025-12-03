package user

import (
	"encoding/json"
	"net/http"

	"github.com/amangirdhar210/meeting-room/internal/adapters/httputil"
	"github.com/amangirdhar210/meeting-room/internal/core/domain"
	"github.com/amangirdhar210/meeting-room/internal/core/service"
	"github.com/amangirdhar210/meeting-room/internal/http/dto"
	"github.com/gorilla/mux"
)

type Handler struct {
	userService service.UserService
}

func NewHandler(userService service.UserService) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	_, role, ok := httputil.GetUserIDRole(r.Context())
	if !ok || role != "admin" {
		httputil.RespondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	var req dto.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	if err := h.userService.Register(user); err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.RespondWithJSON(w, http.StatusCreated, dto.GenericResponse{Message: "user registered successfully"})
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	_, role, ok := httputil.GetUserIDRole(r.Context())
	if !ok || role != "admin" {
		httputil.RespondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	users, err := h.userService.GetAllUsers()
	if err != nil {
		httputil.HandleError(w, err)
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

	httputil.RespondWithJSON(w, http.StatusOK, resp)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, role, ok := httputil.GetUserIDRole(r.Context())
	if !ok || role != "admin" {
		httputil.RespondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		httputil.RespondWithError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	if id == userId {
		httputil.RespondWithError(w, http.StatusForbidden, "cannot delete yourself")
		return
	}
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		httputil.RespondWithError(w, http.StatusNotFound, "user not found")
		return
	}
	if user.Email == "admin@wg.com" {
		httputil.RespondWithError(w, http.StatusForbidden, "cannot delete superadmin")
		return
	}

	if err := h.userService.DeleteUserByID(id); err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.RespondWithJSON(w, http.StatusOK, dto.GenericResponse{Message: "user deleted successfully"})
}
