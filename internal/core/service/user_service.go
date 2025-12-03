package service

import (
	"strings"
	"time"

	"github.com/amangirdhar210/meeting-room/internal/core/domain"
	"github.com/amangirdhar210/meeting-room/internal/core/ports"
	"github.com/google/uuid"
)

type userService struct {
	repo           ports.UserRepository
	passwordHasher ports.PasswordHasher
}

func NewUserService(repo ports.UserRepository, hasher ports.PasswordHasher) UserService {
	return &userService{
		repo:           repo,
		passwordHasher: hasher,
	}
}

func (s *userService) Register(user *domain.User) error {
	if user == nil {
		return domain.ErrInvalidInput
	}

	user.Email = strings.TrimSpace(user.Email)
	user.Name = strings.TrimSpace(user.Name)
	user.Role = strings.TrimSpace(user.Role)
	user.Password = strings.TrimSpace(user.Password)

	if user.Email == "" || user.Password == "" || user.Name == "" || user.Role == "" {
		return domain.ErrInvalidInput
	}

	existing, err := s.repo.FindByEmail(user.Email)
	if err != nil && err != domain.ErrNotFound {
		return err
	}
	if existing != nil {
		return domain.ErrConflict
	}

	hashed, err := s.passwordHasher.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.ID = uuid.New().String()
	user.Password = hashed
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()

	return s.repo.Create(user)
}

func (s *userService) GetAllUsers() ([]domain.User, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, domain.ErrNotFound
	}
	return users, nil
}

func (s *userService) GetUserByID(id string) (*domain.User, error) {
	if id == "" {
		return nil, domain.ErrInvalidInput
	}
	return s.repo.GetByID(id)
}

func (s *userService) DeleteUserByID(id string) error {
	if id == "" {
		return domain.ErrInvalidInput
	}
	return s.repo.DeleteByID(id)
}
