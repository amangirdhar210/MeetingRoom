package app

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/amangirdhar210/meeting-room/internal/http/handlers"
	"github.com/amangirdhar210/meeting-room/internal/http/middleware"
)

func NewRouter(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	roomHandler *handlers.RoomHandler,
	bookingHandler *handlers.BookingHandler,
	jwtSecret string,
) *mux.Router {
	r := mux.NewRouter()

	// Global middlewares
	r.Use(middleware.LoggingMiddleware)

	// Public routes (no JWT)
	auth := r.PathPrefix("/api/auth").Subrouter()
	auth.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)
	auth.HandleFunc("/register", userHandler.RegisterUser).Methods(http.MethodPost)

	// Protected routes (JWT)
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTAuthMiddleware(jwtSecret))

	// Users
	api.HandleFunc("/users", userHandler.GetAllUsers).Methods(http.MethodGet)
	api.HandleFunc("/users/{id:[0-9]+}", userHandler.GetAllUsers).Methods(http.MethodGet)

	// Rooms
	api.HandleFunc("/rooms", roomHandler.GetAllRooms).Methods(http.MethodGet)
	api.HandleFunc("/rooms", roomHandler.AddRoom).Methods(http.MethodPost)
	api.HandleFunc("/rooms/{id:[0-9]+}", roomHandler.GetRoomByID).Methods(http.MethodGet)

	// Bookings
	api.HandleFunc("/bookings", bookingHandler.CreateBooking).Methods(http.MethodPost)
	api.HandleFunc("/bookings", bookingHandler.GetAllBookings).Methods(http.MethodGet)
	api.HandleFunc("/bookings/{id:[0-9]+}/cancel", bookingHandler.CancelBooking).Methods(http.MethodPost)

	return r
}
