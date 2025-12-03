package service

import (
	"strings"

	"github.com/amangirdhar210/meeting-room/internal/core/domain"
	"github.com/amangirdhar210/meeting-room/internal/core/ports"
)

type authService struct {
	userRepo       ports.UserRepository
	tokenGenerator ports.TokenGenerator
	passwordHasher ports.PasswordHasher
}

func NewAuthService(repo ports.UserRepository, tokenGen ports.TokenGenerator, hasher ports.PasswordHasher) AuthService {
	return &authService{
		userRepo:       repo,
		tokenGenerator: tokenGen,
		passwordHasher: hasher,
	}
}

func (s *authService) Login(email, password string) (string, *domain.User, error) {
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

	if !s.passwordHasher.VerifyPassword(user.Password, password) {
		return "", nil, domain.ErrUnauthorized
	}

	token, err := s.tokenGenerator.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
