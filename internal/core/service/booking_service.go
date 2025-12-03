package service

import (
	"time"

	"github.com/amangirdhar210/meeting-room/internal/core/domain"
	"github.com/amangirdhar210/meeting-room/internal/core/ports"
	"github.com/amangirdhar210/meeting-room/internal/pkg/utils"
	"github.com/google/uuid"
)

type bookingService struct {
	repo     ports.BookingRepository
	roomRepo ports.RoomRepository
	userRepo ports.UserRepository
}

func NewBookingService(bRepo ports.BookingRepository, rRepo ports.RoomRepository, uRepo ports.UserRepository) BookingService {
	return &bookingService{
		repo:     bRepo,
		roomRepo: rRepo,
		userRepo: uRepo,
	}
}

func (s *bookingService) CreateBooking(booking *domain.Booking) error {
	if booking == nil {
		return domain.ErrInvalidInput
	}

	if booking.UserID == "" || booking.RoomID == "" {
		return domain.ErrInvalidInput
	}
	if !utils.IsTimeRangeValid(booking.StartTime, booking.EndTime) {
		return domain.ErrTimeRangeInvalid
	}

	user, err := s.userRepo.GetByID(booking.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return domain.ErrNotFound
	}

	room, err := s.roomRepo.GetByID(booking.RoomID)
	if err != nil {
		return err
	}
	if room == nil {
		return domain.ErrNotFound
	}

	existingBookings, err := s.repo.GetByRoomAndTime(booking.RoomID, booking.StartTime, booking.EndTime)
	if err != nil {
		return err
	}

	for _, b := range existingBookings {
		if utils.Overlaps(booking.StartTime, booking.EndTime, b.StartTime, b.EndTime) {
			return domain.ErrRoomUnavailable
		}
	}

	booking.ID = uuid.New().String()
	booking.Status = "confirmed"
	booking.CreatedAt = time.Now().Unix()
	booking.UpdatedAt = time.Now().Unix()

	err = s.repo.Create(booking)
	if err != nil {
		return err
	}

	return nil
}

func (s *bookingService) GetBookingByID(bookingID string) (*domain.Booking, error) {
	if bookingID == "" {
		return nil, domain.ErrInvalidInput
	}

	booking, err := s.repo.GetByID(bookingID)
	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, domain.ErrNotFound
	}

	return booking, nil
}

func (s *bookingService) CancelBooking(bookingID string) error {
	if bookingID == "" {
		return domain.ErrInvalidInput
	}

	booking, err := s.repo.GetByID(bookingID)
	if err != nil {
		return err
	}
	if booking == nil {
		return domain.ErrNotFound
	}

	if booking.Status == "cancelled" {
		return domain.ErrConflict
	}

	err = s.repo.Cancel(bookingID)
	if err != nil {
		return err
	}

	return nil
}

func (s *bookingService) GetAllBookings() ([]domain.Booking, error) {
	bookings, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	if len(bookings) == 0 {
		return nil, domain.ErrNotFound
	}
	return bookings, nil
}

func (s *bookingService) GetBookingsByRoomID(roomID string) ([]domain.Booking, error) {
	if roomID == "" {
		return nil, domain.ErrInvalidInput
	}

	bookings, err := s.repo.GetByRoomID(roomID)
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *bookingService) GetBookingsWithDetailsByRoomID(roomID string) ([]domain.BookingWithDetails, error) {
	if roomID == "" {
		return nil, domain.ErrInvalidInput
	}

	bookings, err := s.repo.GetByRoomID(roomID)
	if err != nil {
		return nil, err
	}

	room, err := s.roomRepo.GetByID(roomID)
	if err != nil {
		return nil, err
	}

	var detailedBookings []domain.BookingWithDetails
	for _, booking := range bookings {
		user, err := s.userRepo.GetByID(booking.UserID)
		if err != nil {
			continue
		}

		detailedBookings = append(detailedBookings, domain.BookingWithDetails{
			Booking:    booking,
			UserName:   user.Name,
			UserEmail:  user.Email,
			RoomName:   room.Name,
			RoomNumber: room.RoomNumber,
		})
	}

	return detailedBookings, nil
}

func (s *bookingService) GetBookingsByUserID(userID string) ([]domain.Booking, error) {
	if userID == "" {
		return nil, domain.ErrInvalidInput
	}

	bookings, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *bookingService) GetBookingsByDateRange(startDate, endDate int64) ([]domain.Booking, error) {
	bookings, err := s.repo.GetByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *bookingService) GetRoomScheduleByDate(roomID string, targetDate int64) (*domain.RoomScheduleResponse, error) {
	if roomID == "" {
		return nil, domain.ErrInvalidInput
	}

	room, err := s.roomRepo.GetByID(roomID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, domain.ErrNotFound
	}

	bookings, err := s.repo.GetByRoomIDAndDate(roomID, targetDate)
	if err != nil {
		return nil, err
	}

	var scheduleSlots []domain.ScheduleSlot
	for _, booking := range bookings {
		user, err := s.userRepo.GetByID(booking.UserID)
		userName := ""
		if err == nil && user != nil {
			userName = user.Name
		}

		scheduleSlots = append(scheduleSlots, domain.ScheduleSlot{
			StartTime: time.Unix(booking.StartTime, 0).Format(time.RFC3339),
			EndTime:   time.Unix(booking.EndTime, 0).Format(time.RFC3339),
			IsBooked:  true,
			BookingID: &booking.ID,
			UserName:  userName,
			Purpose:   booking.Purpose,
		})
	}

	response := &domain.RoomScheduleResponse{
		RoomID:     room.ID,
		RoomName:   room.Name,
		RoomNumber: room.RoomNumber,
		Date:       time.Unix(targetDate, 0).Format("2006-01-02"),
		Bookings:   scheduleSlots,
	}

	return response, nil
}
