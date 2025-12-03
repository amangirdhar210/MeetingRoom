package auth

import (
	"github.com/amangirdhar210/meeting-room/internal/core/domain"
	"golang.org/x/crypto/bcrypt"
)

type bcryptHasher struct{}

func NewBcryptHasher() *bcryptHasher {
	return &bcryptHasher{}
}

func (b *bcryptHasher) HashPassword(password string) (string, error) {
	if password == "" {
		return "", domain.ErrInvalidInput
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (b *bcryptHasher) VerifyPassword(hashed, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}
