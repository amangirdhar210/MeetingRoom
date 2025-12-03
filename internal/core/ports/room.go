package ports

import "github.com/amangirdhar210/meeting-room/internal/core/domain"

type RoomRepository interface {
	Create(room *domain.Room) error
	GetAll() ([]domain.Room, error)
	GetByID(id string) (*domain.Room, error)
	UpdateAvailability(id string, status string) error
	DeleteByID(id string) error
	SearchWithFilters(minCapacity, maxCapacity int, floor *int, amenities []string) ([]domain.Room, error)
}
