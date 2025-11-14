package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/amangirdhar210/meeting-room/internal/domain"
	"github.com/amangirdhar210/meeting-room/internal/http/dto"
	"github.com/amangirdhar210/meeting-room/internal/http/middleware"
	"github.com/gorilla/mux"
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
	userID, role, ok := middleware.GetUserIDRole(r.Context())
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var bookings []domain.Booking
	var err error

	if role == "admin" {
		bookings, err = h.BookingService.GetAllBookings()
	} else {
		bookings, err = h.BookingService.GetBookingsByUserID(userID)
	}

	if err != nil {
		if err == domain.ErrNotFound {
			json.NewEncoder(w).Encode([]dto.BookingDTO{})
			return
		}
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	userID, role, ok := middleware.GetUserIDRole(r.Context())
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	bookingID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid booking id"}`, http.StatusBadRequest)
		return
	}

	if role != "admin" {
		booking, err := h.BookingService.GetBookingByID(bookingID)
		if err != nil {
			if err == domain.ErrNotFound {
				http.Error(w, `{"error":"booking not found"}`, http.StatusNotFound)
			} else {
				http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
			}
			return
		}
		if booking.UserID != userID {
			http.Error(w, `{"error":"forbidden: you can only cancel your own bookings"}`, http.StatusForbidden)
			return
		}
	}

	if err := h.BookingService.CancelBooking(bookingID); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(dto.GenericResponse{Message: "booking canceled successfully"})
}

func (h *BookingHandler) GetMyBookings(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := middleware.GetUserIDRole(r.Context())
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	bookings, err := h.BookingService.GetBookingsByUserID(userID)
	if err != nil {
		if err == domain.ErrNotFound {
			json.NewEncoder(w).Encode([]dto.BookingDTO{})
			return
		}
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *BookingHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomIDStr := vars["id"]
	roomID, err := strconv.ParseInt(roomIDStr, 10, 64)
	if err != nil || roomID <= 0 {
		http.Error(w, `{"error":"invalid room id"}`, http.StatusBadRequest)
		return
	}

	detailedBookings, err := h.BookingService.GetBookingsWithDetailsByRoomID(roomID)
	if err != nil {
		if err == domain.ErrNotFound {
			json.NewEncoder(w).Encode([]dto.DetailedBookingDTO{})
			return
		}
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	var response []dto.DetailedBookingDTO
	for _, booking := range detailedBookings {
		durationMinutes := int(booking.EndTime.Sub(booking.StartTime).Minutes())
		response = append(response, dto.DetailedBookingDTO{
			ID:         booking.ID,
			UserID:     booking.UserID,
			UserName:   booking.UserName,
			UserEmail:  booking.UserEmail,
			RoomID:     booking.RoomID,
			RoomName:   booking.RoomName,
			RoomNumber: booking.RoomNumber,
			StartTime:  booking.StartTime,
			EndTime:    booking.EndTime,
			Duration:   durationMinutes,
			Purpose:    booking.Purpose,
			Status:     booking.Status,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
