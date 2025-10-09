package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/amangirdhar210/meeting-room/internal/domain"
	"github.com/amangirdhar210/meeting-room/internal/http/dto"
	"github.com/amangirdhar210/meeting-room/internal/http/middleware"
)

type RoomHandler struct {
	RoomService domain.RoomService
}

// AddRoom handles POST /api/rooms
func (h *RoomHandler) AddRoom(w http.ResponseWriter, r *http.Request) {
	_, role, ok := middleware.GetUserIDRole(r.Context())
	if !ok || role != "admin" {
		http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
		return
	}

	var req dto.AddRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	room := &domain.Room{
		Name:        req.Name,
		Capacity:    req.Capacity,
		Location:    req.Location,
		IsAvailable: true,
	}

	if err := h.RoomService.AddRoom(room); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.GenericResponse{Message: "room added successfully"})
}

// GetAllRooms handles GET /api/rooms
func (h *RoomHandler) GetAllRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.RoomService.GetAllRooms()
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	var resp []dto.RoomDTO
	for _, rm := range rooms {
		resp = append(resp, dto.RoomDTO{
			ID:          rm.ID,
			Name:        rm.Name,
			Capacity:    rm.Capacity,
			Location:    rm.Location,
			IsAvailable: rm.IsAvailable,
		})
	}

	json.NewEncoder(w).Encode(resp)
}

// GetRoomByID handles GET /api/rooms/{id}
func (h *RoomHandler) GetRoomByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid room id"}`, http.StatusBadRequest)
		return
	}

	room, err := h.RoomService.GetRoomByID(id)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	resp := dto.RoomDTO{
		ID:          room.ID,
		Name:        room.Name,
		Capacity:    room.Capacity,
		Location:    room.Location,
		IsAvailable: room.IsAvailable,
	}

	json.NewEncoder(w).Encode(resp)
}
