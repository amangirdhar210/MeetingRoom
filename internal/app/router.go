package app

import (
	"database/sql"
	"net/http"

	"github.com/amangirdhar210/meeting-room/internal/http/handlers"
	"github.com/amangirdhar210/meeting-room/internal/http/middleware"
	"github.com/amangirdhar210/meeting-room/internal/repositories/mysql"
	"github.com/amangirdhar210/meeting-room/internal/service"
	"github.com/gorilla/mux"
)

// SetupRouter wires repositories -> services -> handlers and registers routes.
// It intentionally uses direct struct literals for handlers so no constructor
// functions are required in your handlers package.
func SetupRouter(db *sql.DB, jwtSecret string) http.Handler {
	// repositories
	userRepo := mysql.NewUserRepositoryMySQL(db)
	roomRepo := mysql.NewRoomRepositoryMySQL(db)
	bookingRepo := mysql.NewBookingRepositoryMySQL(db)

	// services (AuthService uses userRepo; JWT secret should be initialized in main)
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo, authService)
	roomService := service.NewRoomService(roomRepo)
	bookingService := service.NewBookingService(bookingRepo, roomRepo, userRepo)

	// handlers (use literals so you don't need NewXxx constructors)
	authHandler := &handlers.AuthHandler{AuthService: authService}
	userHandler := &handlers.UserHandler{UserService: userService}
	roomHandler := &handlers.RoomHandler{RoomService: roomService}
	bookingHandler := &handlers.BookingHandler{BookingService: bookingService}

	// router
	r := mux.NewRouter()

	// public
	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// register/login
	r.HandleFunc("/api/register", userHandler.RegisterUser).Methods("POST")
	r.HandleFunc("/api/login", authHandler.Login).Methods("POST")

	// protected API subrouter
	api := r.PathPrefix("/api").Subrouter()
	// logging first, then JWT auth
	api.Use(middleware.LoggingMiddleware)
	api.Use(middleware.JWTAuthMiddleware(jwtSecret))

	// users
	api.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")

	// rooms
	api.HandleFunc("/rooms", roomHandler.AddRoom).Methods("POST")
	api.HandleFunc("/rooms", roomHandler.GetAllRooms).Methods("GET")
	api.HandleFunc("/rooms/{id:[0-9]+}", roomHandler.GetRoomByID).Methods("GET")

	// bookings
	api.HandleFunc("/bookings", bookingHandler.CreateBooking).Methods("POST")
	api.HandleFunc("/bookings", bookingHandler.GetAllBookings).Methods("GET")
	api.HandleFunc("/bookings/{id:[0-9]+}", bookingHandler.CancelBooking).Methods("DELETE")

	return r
}
