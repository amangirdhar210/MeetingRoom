package booking

import (
	"encoding/json"
	"net/http"
	"time"

	httputil "github.com/amangirdhar210/meeting-room/internal/adapters/httpUtils"
	"github.com/amangirdhar210/meeting-room/internal/core/domain"
	"github.com/amangirdhar210/meeting-room/internal/core/service"
	"github.com/amangirdhar210/meeting-room/internal/http/dto"
	"github.com/gorilla/mux"
)

type Handler struct {
	bookingService service.BookingService
}

func NewHandler(bookingService service.BookingService) *Handler {
	return &Handler{bookingService: bookingService}
}

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := httputil.GetUserIDRole(r.Context())
	if !ok {
		httputil.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		httputil.RespondWithError(w, http.StatusBadRequest, "invalid start_time format")
		return
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		httputil.RespondWithError(w, http.StatusBadRequest, "invalid end_time format")
		return
	}

	booking := &domain.Booking{
		UserID:    userID,
		RoomID:    req.RoomID,
		StartTime: startTime.Unix(),
		EndTime:   endTime.Unix(),
		Purpose:   req.Purpose,
	}

	if err := h.bookingService.CreateBooking(booking); err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.RespondWithJSON(w, http.StatusCreated, dto.GenericResponse{Message: "booking created successfully"})
}

func (h *Handler) GetAllBookings(w http.ResponseWriter, r *http.Request) {
	userID, role, ok := httputil.GetUserIDRole(r.Context())
	if !ok {
		httputil.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var bookings []domain.Booking
	var err error

	if role == "admin" {
		bookings, err = h.bookingService.GetAllBookings()
	} else {
		bookings, err = h.bookingService.GetBookingsByUserID(userID)
	}

	if err != nil {
		if err == domain.ErrNotFound {
			httputil.RespondWithJSON(w, http.StatusOK, []dto.BookingDTO{})
			return
		}
		httputil.HandleError(w, err)
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

	httputil.RespondWithJSON(w, http.StatusOK, resp)
}

func (h *Handler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	userID, role, ok := httputil.GetUserIDRole(r.Context())
	if !ok {
		httputil.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	vars := mux.Vars(r)
	bookingID := vars["id"]
	if bookingID == "" {
		httputil.RespondWithError(w, http.StatusBadRequest, "invalid booking id")
		return
	}

	if role != "admin" {
		booking, err := h.bookingService.GetBookingByID(bookingID)
		if err != nil {
			if err == domain.ErrNotFound {
				httputil.RespondWithError(w, http.StatusNotFound, "booking not found")
			} else {
				httputil.HandleError(w, err)
			}
			return
		}
		if booking.UserID != userID {
			httputil.RespondWithError(w, http.StatusForbidden, "forbidden: you can only cancel your own bookings")
			return
		}
	}

	if err := h.bookingService.CancelBooking(bookingID); err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.RespondWithJSON(w, http.StatusOK, dto.GenericResponse{Message: "booking canceled successfully"})
}

func (h *Handler) GetMyBookings(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := httputil.GetUserIDRole(r.Context())
	if !ok {
		httputil.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	bookings, err := h.bookingService.GetBookingsByUserID(userID)
	if err != nil {
		if err == domain.ErrNotFound {
			httputil.RespondWithJSON(w, http.StatusOK, []dto.BookingDTO{})
			return
		}
		httputil.HandleError(w, err)
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

	httputil.RespondWithJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["id"]
	if roomID == "" {
		httputil.RespondWithError(w, http.StatusBadRequest, "invalid room id")
		return
	}

	detailedBookings, err := h.bookingService.GetBookingsWithDetailsByRoomID(roomID)
	if err != nil {
		if err == domain.ErrNotFound {
			httputil.RespondWithJSON(w, http.StatusOK, []dto.DetailedBookingDTO{})
			return
		}
		httputil.HandleError(w, err)
		return
	}

	var response []dto.DetailedBookingDTO
	for _, booking := range detailedBookings {
		durationMinutes := int((booking.EndTime - booking.StartTime) / 60)
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

	httputil.RespondWithJSON(w, http.StatusOK, response)
}

func (h *Handler) GetScheduleByDate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["id"]
	if roomID == "" {
		httputil.RespondWithError(w, http.StatusBadRequest, "invalid room id")
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		httputil.RespondWithError(w, http.StatusBadRequest, "date parameter is required (format: YYYY-MM-DD)")
		return
	}

	targetDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		httputil.RespondWithError(w, http.StatusBadRequest, "invalid date format, use YYYY-MM-DD")
		return
	}

	schedule, err := h.bookingService.GetRoomScheduleByDate(roomID, targetDate.Unix())
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.RespondWithJSON(w, http.StatusOK, schedule)
}
