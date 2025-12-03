package ports

type TokenGenerator interface {
	GenerateToken(userID, role string) (string, error)
}

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashed, plain string) bool
}
