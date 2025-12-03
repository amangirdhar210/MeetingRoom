package main

import (
	"fmt"
	"log"

	"github.com/amangirdhar210/meeting-room/internal/adapters/auth"
	httpAdapter "github.com/amangirdhar210/meeting-room/internal/adapters/http"
	"github.com/amangirdhar210/meeting-room/internal/adapters/repository"
	"github.com/amangirdhar210/meeting-room/internal/config"
	"github.com/amangirdhar210/meeting-room/internal/core/service"
	"github.com/amangirdhar210/meeting-room/internal/repositories/mysql"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it")
	}

	cfg := config.LoadConfig()

	if cfg.JWT.Secret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	dbCfg := mysql.DBConfig{
		Path: cfg.Database.Path,
	}

	db, err := mysql.NewSQLiteConnection(dbCfg)
	if err != nil {
		log.Fatalf("Failed to connect to SQLite: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	if err := mysql.InitSQLite(db); err != nil {
		log.Fatalf("Failed to initialize SQLite schema: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	jwtGenerator := auth.NewJWTGenerator(cfg.JWT.Secret, cfg.JWT.ExpirationTime)
	passwordHasher := auth.NewBcryptHasher()

	authService := service.NewAuthService(userRepo, jwtGenerator, passwordHasher)
	userService := service.NewUserService(userRepo, passwordHasher)
	roomService := service.NewRoomService(roomRepo)
	bookingService := service.NewBookingService(bookingRepo, roomRepo, userRepo)

	server := httpAdapter.NewHTTPServer(
		cfg,
		userService,
		authService,
		roomService,
		bookingService,
		jwtGenerator,
	)

	fmt.Printf("Server starting on http://localhost%s\n", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
