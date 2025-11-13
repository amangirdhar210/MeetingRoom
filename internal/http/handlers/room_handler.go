package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/amangirdhar210/meeting-room/internal/domain"
	"github.com/amangirdhar210/meeting-room/internal/http/dto"
	"github.com/amangirdhar210/meeting-room/internal/http/middleware"
	"github.com/gorilla/mux"
)

type RoomHandler struct {
	RoomService domain.RoomService
}

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
		RoomNumber:  req.RoomNumber,
		Capacity:    req.Capacity,
		Floor:       req.Floor,
		Amenities:   req.Amenities,
		Status:      req.Status,
		Location:    req.Location,
		Description: req.Description,
	}

	if room.Status == "" {
		room.Status = "Available"
	}
	if room.Amenities == nil {
		room.Amenities = []string{}
	}

	if err := h.RoomService.AddRoom(room); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.GenericResponse{Message: "room added successfully"})
}

func (h *RoomHandler) GetAllRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.RoomService.GetAllRooms()
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	var resp []dto.RoomDTO
	for _, rm := range rooms {
		resp = append(resp, dto.RoomDTO{
			ID:          rm.ID,
			Name:        rm.Name,
			RoomNumber:  rm.RoomNumber,
			Capacity:    rm.Capacity,
			Floor:       rm.Floor,
			Amenities:   rm.Amenities,
			Status:      rm.Status,
			Location:    rm.Location,
			Description: rm.Description,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *RoomHandler) GetRoomByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid room id"}`, http.StatusBadRequest)
		return
	}

	room, err := h.RoomService.GetRoomByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, `{"error":"resource not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		}
		return
	}

	resp := dto.RoomDTO{
		ID:          room.ID,
		Name:        room.Name,
		RoomNumber:  room.RoomNumber,
		Capacity:    room.Capacity,
		Floor:       room.Floor,
		Amenities:   room.Amenities,
		Status:      room.Status,
		Location:    room.Location,
		Description: room.Description,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *RoomHandler) DeleteRoomByID(w http.ResponseWriter, r *http.Request) {
	_, role, ok := middleware.GetUserIDRole(r.Context())
	if !ok || role != "admin" {
		http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid room id"}`, http.StatusBadRequest)
		return
	}

	if err := h.RoomService.DeleteRoomByID(id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, `{"error":"resource not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.GenericResponse{Message: "room deleted successfully"})
}
