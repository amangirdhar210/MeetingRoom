package http

import (
	"net/http"

	"github.com/amangirdhar210/meeting-room/internal/adapters/auth"
	authHandler "github.com/amangirdhar210/meeting-room/internal/adapters/http/auth"
	bookingHandler "github.com/amangirdhar210/meeting-room/internal/adapters/http/booking"
	roomHandler "github.com/amangirdhar210/meeting-room/internal/adapters/http/room"
	userHandler "github.com/amangirdhar210/meeting-room/internal/adapters/http/user"
	"github.com/amangirdhar210/meeting-room/internal/config"
	"github.com/amangirdhar210/meeting-room/internal/core/service"
	"github.com/gorilla/mux"
)

func NewHTTPServer(cfg *config.Config, userService service.UserService, authService service.AuthService, roomService service.RoomService, bookingService service.BookingService, jwtGenerator *auth.JWTGenerator) *http.Server {
	authH := authHandler.NewHandler(authService)
	userH := userHandler.NewHandler(userService)
	roomH := roomHandler.NewHandler(roomService)
	bookingH := bookingHandler.NewHandler(bookingService)

	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	}).Methods("GET")

	router.HandleFunc("/api/login", authH.Login).Methods("POST")

	api := router.PathPrefix("/api").Subrouter()
	api.Use(LoggingMiddleware)
	api.Use(JWTAuthMiddleware(jwtGenerator))

	api.HandleFunc("/users", userH.GetAllUsers).Methods("GET")
	api.HandleFunc("/users/{id}", userH.DeleteUser).Methods("DELETE")
	api.HandleFunc("/register", userH.RegisterUser).Methods("POST")

	api.HandleFunc("/rooms", roomH.AddRoom).Methods("POST")
	api.HandleFunc("/rooms", roomH.GetAllRooms).Methods("GET")
	api.HandleFunc("/rooms/search", roomH.SearchRooms).Methods("GET")
	api.HandleFunc("/rooms/check-availability", roomH.CheckAvailability).Methods("POST")
	api.HandleFunc("/rooms/{id}", roomH.GetRoomByID).Methods("GET")
	api.HandleFunc("/rooms/{id}/delete", roomH.DeleteRoomByID).Methods("DELETE")
	api.HandleFunc("/rooms/{id}/schedule", bookingH.GetSchedule).Methods("GET")
	api.HandleFunc("/rooms/{id}/schedule/date", bookingH.GetScheduleByDate).Methods("GET")

	api.HandleFunc("/bookings", bookingH.CreateBooking).Methods("POST")
	api.HandleFunc("/bookings", bookingH.GetAllBookings).Methods("GET")
	api.HandleFunc("/bookings/my", bookingH.GetMyBookings).Methods("GET")
	api.HandleFunc("/bookings/{id}", bookingH.CancelBooking).Methods("DELETE")

	wrappedRouter := CORSMiddleware(router, cfg.CORS.AllowedOrigins)

	server := &http.Server{
		Handler:      wrappedRouter,
		Addr:         cfg.Server.Port,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return server
}
