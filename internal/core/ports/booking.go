package ports

import "github.com/amangirdhar210/meeting-room/internal/core/domain"

type BookingRepository interface {
	Create(booking *domain.Booking) error
	GetByID(id string) (*domain.Booking, error)
	GetAll() ([]domain.Booking, error)
	GetByRoomAndTime(roomID string, start, end int64) ([]domain.Booking, error)
	GetByRoomID(roomID string) ([]domain.Booking, error)
	GetByUserID(userID string) ([]domain.Booking, error)
	Cancel(id string) error
	GetByDateRange(startDate, endDate int64) ([]domain.Booking, error)
	GetByRoomIDAndDate(roomID string, date int64) ([]domain.Booking, error)
}
