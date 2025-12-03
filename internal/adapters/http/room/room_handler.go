package room

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/amangirdhar210/meeting-room/internal/adapters/httputil"
	"github.com/amangirdhar210/meeting-room/internal/core/domain"
	"github.com/amangirdhar210/meeting-room/internal/core/service"
	"github.com/amangirdhar210/meeting-room/internal/http/dto"
	"github.com/gorilla/mux"
)

type Handler struct {
	roomService service.RoomService
}

func NewHandler(roomService service.RoomService) *Handler {
	return &Handler{roomService: roomService}
}

func (h *Handler) AddRoom(w http.ResponseWriter, r *http.Request) {
	_, role, ok := httputil.GetUserIDRole(r.Context())
	if !ok || role != "admin" {
		httputil.RespondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	var req dto.AddRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.RespondWithError(w, http.StatusBadRequest, "invalid request body")
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

	if err := h.roomService.AddRoom(room); err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.RespondWithJSON(w, http.StatusCreated, dto.GenericResponse{Message: "room added successfully"})
}

func (h *Handler) GetAllRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.roomService.GetAllRooms()
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		httputil.HandleError(w, err)
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

	httputil.RespondWithJSON(w, http.StatusOK, response)
}

func (h *Handler) GetRoomByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["id"]
	if roomID == "" {
		httputil.RespondWithError(w, http.StatusBadRequest, "invalid room id")
		return
	}

	room, err := h.roomService.GetRoomByID(roomID)
	if err != nil {
		httputil.HandleError(w, err)
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

	httputil.RespondWithJSON(w, http.StatusOK, response)
}

func (h *Handler) DeleteRoomByID(w http.ResponseWriter, r *http.Request) {
	_, userRole, ok := httputil.GetUserIDRole(r.Context())
	if !ok || userRole != "admin" {
		httputil.RespondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	vars := mux.Vars(r)
	roomID := vars["id"]
	if roomID == "" {
		httputil.RespondWithError(w, http.StatusBadRequest, "invalid room id")
		return
	}

	if err := h.roomService.DeleteRoomByID(roomID); err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.RespondWithJSON(w, http.StatusOK, dto.GenericResponse{Message: "room deleted successfully"})
}

func (h *Handler) SearchRooms(w http.ResponseWriter, r *http.Request) {
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

	var startTime, endTime *int64
	if startStr := queryParams.Get("startTime"); startStr != "" {
		if t, err := time.Parse(time.RFC3339, startStr); err == nil {
			unix := t.Unix()
			startTime = &unix
		}
	}
	if endStr := queryParams.Get("endTime"); endStr != "" {
		if t, err := time.Parse(time.RFC3339, endStr); err == nil {
			unix := t.Unix()
			endTime = &unix
		}
	}

	rooms, err := h.roomService.SearchRooms(minCapacity, maxCapacity, floor, amenities, startTime, endTime)
	if err != nil {
		httputil.HandleError(w, err)
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

	httputil.RespondWithJSON(w, http.StatusOK, response)
}

func (h *Handler) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	var request dto.AvailabilityCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		httputil.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	startTime, err := time.Parse(time.RFC3339, request.StartTime)
	if err != nil {
		httputil.RespondWithError(w, http.StatusBadRequest, "invalid start_time format")
		return
	}

	endTime, err := time.Parse(time.RFC3339, request.EndTime)
	if err != nil {
		httputil.RespondWithError(w, http.StatusBadRequest, "invalid end_time format")
		return
	}

	isAvailable, conflictingBookings, err := h.roomService.CheckAvailability(request.RoomID, startTime.Unix(), endTime.Unix())
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	room, err := h.roomService.GetRoomByID(request.RoomID)
	if err != nil {
		httputil.RespondWithError(w, http.StatusNotFound, "room not found")
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
		RequestedStart:   startTime.Unix(),
		RequestedEnd:     endTime.Unix(),
		ConflictingSlots: conflictingSlots,
		SuggestedSlots:   []dto.TimeSlotDTO{},
	}

	httputil.RespondWithJSON(w, http.StatusOK, response)
}
