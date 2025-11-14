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

func SetupRouter(db *sql.DB, jwtSecret string) http.Handler {
	userRepo := mysql.NewUserRepositorySQLite(db)
	roomRepo := mysql.NewRoomRepositorySQLite(db)
	bookingRepo := mysql.NewBookingRepositorySQLite(db)

	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo, authService)
	roomService := service.NewRoomService(roomRepo)
	bookingService := service.NewBookingService(bookingRepo, roomRepo, userRepo)

	authHandler := &handlers.AuthHandler{AuthService: authService}
	userHandler := &handlers.UserHandler{UserService: userService}
	roomHandler := &handlers.RoomHandler{RoomService: roomService}
	bookingHandler := &handlers.BookingHandler{BookingService: bookingService}

	r := mux.NewRouter()

	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	r.HandleFunc("/api/login", authHandler.Login).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.LoggingMiddleware)
	api.Use(middleware.JWTAuthMiddleware(jwtSecret))

	api.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	api.HandleFunc("/users/{id:[0-9]+}", userHandler.DeleteUser).Methods("DELETE")

	api.HandleFunc("/register", userHandler.RegisterUser).Methods("POST")

	api.HandleFunc("/rooms", roomHandler.AddRoom).Methods("POST")
	api.HandleFunc("/rooms", roomHandler.GetAllRooms).Methods("GET")
	api.HandleFunc("/rooms/search", roomHandler.SearchRooms).Methods("GET")
	api.HandleFunc("/rooms/check-availability", roomHandler.CheckAvailability).Methods("POST")
	api.HandleFunc("/rooms/{id:[0-9]+}", roomHandler.GetRoomByID).Methods("GET")
	api.HandleFunc("/rooms/{id}", roomHandler.DeleteRoomByID).Methods("DELETE")
	api.HandleFunc("/rooms/{id:[0-9]+}/schedule", bookingHandler.GetSchedule).Methods("GET")

	api.HandleFunc("/bookings", bookingHandler.CreateBooking).Methods("POST")
	api.HandleFunc("/bookings", bookingHandler.GetAllBookings).Methods("GET")
	api.HandleFunc("/bookings/{id:[0-9]+}", bookingHandler.CancelBooking).Methods("DELETE")

	return r
}
