package ports

import "github.com/amangirdhar210/meeting-room/internal/core/domain"

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	GetByID(id string) (*domain.User, error)
	GetAll() ([]domain.User, error)
	DeleteByID(id string) error
}
