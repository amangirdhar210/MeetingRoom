package service

import (
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/amangirdhar210/meeting-room/internal/domain"
	"github.com/amangirdhar210/meeting-room/internal/pkg/jwt"
)

type AuthService struct {
	userRepo domain.UserRepository
}

func NewAuthService(repo domain.UserRepository) *AuthService {
	return &AuthService{userRepo: repo}
}

func (s *AuthService) HashPassword(password string) (string, error) {
	if password == "" {
		return "", domain.ErrInvalidInput
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *AuthService) VerifyPassword(hashed, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}

func (s *AuthService) AuthenticateUser(email, password string) (string, *domain.User, error) {
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	if email == "" || password == "" {
		return "", nil, domain.ErrInvalidInput
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, domain.ErrUnauthorized
	}

	if !s.VerifyPassword(user.Password, password) {
		return "", nil, domain.ErrUnauthorized
	}

	token, err := jwt.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
