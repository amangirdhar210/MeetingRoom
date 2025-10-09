package service

import (
	"time"

	"github.com/amangirdhar210/meeting-room/internal/domain"
	"github.com/amangirdhar210/meeting-room/internal/pkg/utils"
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

	if booking.UserID <= 0 || booking.RoomID <= 0 {
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

	booking.Status = "confirmed"
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()

	return s.repo.Create(booking)
}

func (s *BookingServiceImpl) CancelBooking(id int64) error {
	if id <= 0 {
		return domain.ErrInvalidInput
	}

	booking, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if booking == nil {
		return domain.ErrNotFound
	}

	if booking.Status == "cancelled" {
		return domain.ErrConflict
	}

	err = s.repo.Cancel(id)
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
