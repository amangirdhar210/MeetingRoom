package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/amangirdhar210/meeting-room/internal/app"
	"github.com/amangirdhar210/meeting-room/internal/http/handlers"
	"github.com/amangirdhar210/meeting-room/internal/repositories/mysql"
	"github.com/amangirdhar210/meeting-room/internal/service"
)

func main() {
	// 1️⃣ Connect to DB (replace DSN with your own)
	dsn := "root:password@tcp(127.0.0.1:3306)/meeting_room?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ failed to connect to DB: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("❌ failed to ping DB: %v", err)
	}

	// 2️⃣ Initialize Repositories
	userRepo := mysql.NewUserRepository(db)
	roomRepo := mysql.NewRoomRepository(db)
	bookingRepo := mysql.NewBookingRepository(db)

	// 3️⃣ Initialize Services
	userService := service.NewUserService(userRepo)
	roomService := service.NewRoomService(roomRepo)
	bookingService := service.NewBookingService(bookingRepo, roomRepo, userRepo)
	authService := service.NewAuthService(userRepo) // if you have Auth service

	// 4️⃣ Initialize Handlers
	userHandler := handlers.NewUserHandler(userService)
	roomHandler := handlers.NewRoomHandler(roomService)
	bookingHandler := handlers.NewBookingHandler(bookingService)
	authHandler := handlers.NewAuthHandler(authService)

	// 5️⃣ Initialize Router
	jwtSecret := "your-secret-key"
	router := app.NewRouter(authHandler, userHandler, roomHandler, bookingHandler, jwtSecret)

	// 6️⃣ Start Server
	app.StartServer(":8080", router)
}
