package service

import (
	"strings"
	"time"

	"github.com/amangirdhar210/meeting-room/internal/domain"
)

type UserServiceImpl struct {
	repo        domain.UserRepository
	authService *AuthService
}

func NewUserService(repo domain.UserRepository, auth *AuthService) domain.UserService {
	return &UserServiceImpl{
		repo:        repo,
		authService: auth,
	}
}

func (s *UserServiceImpl) Register(user *domain.User) error {
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

	existing, _ := s.repo.FindByEmail(user.Email)
	if existing != nil {
		return domain.ErrConflict
	}

	hashed, err := s.authService.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return s.repo.Create(user)
}

func (s *UserServiceImpl) GetAllUsers() ([]domain.User, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, domain.ErrNotFound
	}
	return users, nil
}
