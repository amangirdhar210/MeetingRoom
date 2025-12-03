package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTGenerator struct {
	secret         []byte
	expirationTime time.Duration
}

func NewJWTGenerator(secret string, expirationTime time.Duration) *JWTGenerator {
	return &JWTGenerator{
		secret:         []byte(secret),
		expirationTime: expirationTime,
	}
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (j *JWTGenerator) GenerateToken(userID string, role string) (string, error) {
	expirationTime := time.Now().Add(j.expirationTime)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTGenerator) ValidateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		return j.secret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
