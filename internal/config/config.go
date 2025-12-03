package config

import (
	"os"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	Path            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type JWTConfig struct {
	Secret         string
	ExpirationTime time.Duration
}

type CORSConfig struct {
	AllowedOrigins []string
}

func LoadConfig() *Config {
	jwtSecret := os.Getenv("JWT_SECRET")

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./meeting_room_db.sqlite"
	}

	return &Config{
		Server: ServerConfig{
			Port:            ":8080",
			ReadTimeout:     15 * time.Second,
			WriteTimeout:    15 * time.Second,
			ShutdownTimeout: 30 * time.Second,
		},
		Database: DatabaseConfig{
			Path:            dbPath,
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 5 * time.Minute,
		},
		JWT: JWTConfig{
			Secret:         jwtSecret,
			ExpirationTime: 24 * time.Hour,
		},
		CORS: CORSConfig{
			AllowedOrigins: []string{
				"http://localhost:4200",
				"http://127.0.0.1:4200",
			},
		},
	}
}
