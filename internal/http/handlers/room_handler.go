package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

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

	var response []dto.RoomDTO
	for _, room := range rooms {
		response = append(response, dto.RoomDTO{
			ID:          room.ID,
			Name:        room.Name,
			RoomNumber:  room.RoomNumber,
			Capacity:    room.Capacity,
			Floor:       room.Floor,
			Amenities:   room.Amenities,
			Status:      room.Status,
			Location:    room.Location,
			Description: room.Description,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *RoomHandler) GetRoomByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomIDStr := vars["id"]
	roomID, err := strconv.ParseInt(roomIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid room id"}`, http.StatusBadRequest)
		return
	}

	room, err := h.RoomService.GetRoomByID(roomID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, `{"error":"resource not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		}
		return
	}

	response := dto.RoomDTO{
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
	json.NewEncoder(w).Encode(response)
}

func (h *RoomHandler) DeleteRoomByID(w http.ResponseWriter, r *http.Request) {
	_, userRole, ok := middleware.GetUserIDRole(r.Context())
	if !ok || userRole != "admin" {
		http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	roomIDStr := vars["id"]
	roomID, err := strconv.ParseInt(roomIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid room id"}`, http.StatusBadRequest)
		return
	}

	if err := h.RoomService.DeleteRoomByID(roomID); err != nil {
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

func (h *RoomHandler) SearchRooms(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	minCapacity := 0
	if minCapStr := queryParams.Get("minCapacity"); minCapStr != "" {
		if val, err := strconv.Atoi(minCapStr); err == nil {
			minCapacity = val
		}
	}

	maxCapacity := 0
	if maxCapStr := queryParams.Get("maxCapacity"); maxCapStr != "" {
		if val, err := strconv.Atoi(maxCapStr); err == nil {
			maxCapacity = val
		}
	}

	var floor *int
	if floorStr := queryParams.Get("floor"); floorStr != "" {
		if val, err := strconv.Atoi(floorStr); err == nil {
			floor = &val
		}
	}

	var amenities []string
	if amenitiesStr := queryParams.Get("amenities"); amenitiesStr != "" {
		amenities = []string{}
		if err := json.Unmarshal([]byte(amenitiesStr), &amenities); err == nil {
		} else {
			amenities = []string{amenitiesStr}
		}
	}

	var startTime, endTime *time.Time
	if startStr := queryParams.Get("startTime"); startStr != "" {
		if t, err := time.Parse(time.RFC3339, startStr); err == nil {
			startTime = &t
		}
	}
	if endStr := queryParams.Get("endTime"); endStr != "" {
		if t, err := time.Parse(time.RFC3339, endStr); err == nil {
			endTime = &t
		}
	}

	rooms, err := h.RoomService.SearchRooms(minCapacity, maxCapacity, floor, amenities, startTime, endTime)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	var response []dto.RoomDTO
	for _, room := range rooms {
		response = append(response, dto.RoomDTO{
			ID:          room.ID,
			Name:        room.Name,
			RoomNumber:  room.RoomNumber,
			Capacity:    room.Capacity,
			Floor:       room.Floor,
			Amenities:   room.Amenities,
			Status:      room.Status,
			Location:    room.Location,
			Description: room.Description,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *RoomHandler) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	var request dto.AvailabilityCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse(time.RFC3339, request.StartTime)
	if err != nil {
		http.Error(w, `{"error":"invalid start_time format"}`, http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(time.RFC3339, request.EndTime)
	if err != nil {
		http.Error(w, `{"error":"invalid end_time format"}`, http.StatusBadRequest)
		return
	}

	isAvailable, conflictingBookings, err := h.RoomService.CheckAvailability(request.RoomID, startTime, endTime)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, `{"error":"room not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		}
		return
	}

	room, err := h.RoomService.GetRoomByID(request.RoomID)
	if err != nil {
		http.Error(w, `{"error":"room not found"}`, http.StatusNotFound)
		return
	}

	var conflictingSlots []dto.ConflictingBookingDTO
	for _, conflictBooking := range conflictingBookings {
		conflictingSlots = append(conflictingSlots, dto.ConflictingBookingDTO{
			BookingID: conflictBooking.ID,
			StartTime: conflictBooking.StartTime,
			EndTime:   conflictBooking.EndTime,
			Purpose:   conflictBooking.Purpose,
		})
	}

	response := dto.AvailabilityCheckResponse{
		Available:        isAvailable,
		RoomID:           request.RoomID,
		RoomName:         room.Name,
		RequestedStart:   startTime,
		RequestedEnd:     endTime,
		ConflictingSlots: conflictingSlots,
		SuggestedSlots:   []dto.TimeSlotDTO{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
