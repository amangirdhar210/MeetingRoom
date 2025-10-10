package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/amangirdhar210/meeting-room/internal/domain"
	"github.com/amangirdhar210/meeting-room/internal/http/dto"
	"github.com/amangirdhar210/meeting-room/internal/http/middleware"
)

type BookingHandler struct {
	BookingService domain.BookingService
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := middleware.GetUserIDRole(r.Context())
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req dto.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		http.Error(w, `{"error":"invalid start_time format"}`, http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		http.Error(w, `{"error":"invalid end_time format"}`, http.StatusBadRequest)
		return
	}

	booking := &domain.Booking{
		UserID:    userID,
		RoomID:    req.RoomID,
		StartTime: startTime,
		EndTime:   endTime,
		Purpose:   req.Purpose,
	}

	if err := h.BookingService.CreateBooking(booking); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.GenericResponse{Message: "booking created successfully"})
}

func (h *BookingHandler) GetAllBookings(w http.ResponseWriter, r *http.Request) {
	_, role, ok := middleware.GetUserIDRole(r.Context())
	if !ok || role != "admin" {
		http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
		return
	}

	bookings, err := h.BookingService.GetAllBookings()
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	var resp []dto.BookingDTO
	for _, b := range bookings {
		resp = append(resp, dto.BookingDTO{
			ID:        b.ID,
			UserID:    b.UserID,
			RoomID:    b.RoomID,
			StartTime: b.StartTime,
			EndTime:   b.EndTime,
			Purpose:   b.Purpose,
		})
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	_, role, ok := middleware.GetUserIDRole(r.Context())
	if !ok || role != "admin" {
		http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid booking id"}`, http.StatusBadRequest)
		return
	}

	if err := h.BookingService.CancelBooking(id); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(dto.GenericResponse{Message: "booking canceled successfully"})
}
