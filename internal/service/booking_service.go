package service

import (
	"time"

	"github.com/amangirdhar210/meeting-room/internal/domain"
	"github.com/amangirdhar210/meeting-room/internal/pkg/utils"
	"github.com/google/uuid"
)

type BookingServiceImpl struct {
	repo     domain.BookingRepository
	roomRepo domain.RoomRepository
	userRepo domain.UserRepository
}

func NewBookingService(bRepo domain.BookingRepository, rRepo domain.RoomRepository, uRepo domain.UserRepository) domain.BookingService {
	return &BookingServiceImpl{
		repo:     bRepo,
		roomRepo: rRepo,
		userRepo: uRepo,
	}
}

func (s *BookingServiceImpl) CreateBooking(booking *domain.Booking) error {
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
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()

	err = s.repo.Create(booking)
	if err != nil {
		return err
	}

	return nil
}

func (s *BookingServiceImpl) GetBookingByID(bookingID string) (*domain.Booking, error) {
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

func (s *BookingServiceImpl) CancelBooking(bookingID string) error {
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

func (s *BookingServiceImpl) GetAllBookings() ([]domain.Booking, error) {
	bookings, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	if len(bookings) == 0 {
		return nil, domain.ErrNotFound
	}
	return bookings, nil
}

func (s *BookingServiceImpl) GetBookingsByRoomID(roomID string) ([]domain.Booking, error) {
	if roomID == "" {
		return nil, domain.ErrInvalidInput
	}

	bookings, err := s.repo.GetByRoomID(roomID)
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *BookingServiceImpl) GetBookingsWithDetailsByRoomID(roomID string) ([]domain.BookingWithDetails, error) {
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

func (s *BookingServiceImpl) GetBookingsByUserID(userID string) ([]domain.Booking, error) {
	if userID == "" {
		return nil, domain.ErrInvalidInput
	}

	bookings, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *BookingServiceImpl) GetBookingsByDateRange(startDate, endDate time.Time) ([]domain.Booking, error) {
	bookings, err := s.repo.GetByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *BookingServiceImpl) GetRoomScheduleByDate(roomID string, targetDate time.Time) (*domain.RoomScheduleResponse, error) {
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
			StartTime: booking.StartTime.Format(time.RFC3339),
			EndTime:   booking.EndTime.Format(time.RFC3339),
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
		Date:       targetDate.Format("2006-01-02"),
		Bookings:   scheduleSlots,
	}

	return response, nil
}
