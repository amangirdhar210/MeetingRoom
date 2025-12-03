package service

import "github.com/amangirdhar210/meeting-room/internal/core/domain"

type UserService interface {
	Register(user *domain.User) error
	GetAllUsers() ([]domain.User, error)
	GetUserByID(id string) (*domain.User, error)
	DeleteUserByID(id string) error
}

type AuthService interface {
	Login(email, password string) (token string, user *domain.User, err error)
}

type RoomService interface {
	AddRoom(room *domain.Room) error
	GetAllRooms() ([]domain.Room, error)
	GetRoomByID(id string) (*domain.Room, error)
	DeleteRoomByID(id string) error
	SearchRooms(minCapacity, maxCapacity int, floor *int, amenities []string, startTime, endTime *int64) ([]domain.Room, error)
	CheckAvailability(roomID string, startTime, endTime int64) (bool, []domain.Booking, error)
	GetAvailableSlots(roomID string, date int64, slotDuration int) ([]domain.TimeSlot, error)
}

type BookingService interface {
	CreateBooking(booking *domain.Booking) error
	GetBookingByID(bookingID string) (*domain.Booking, error)
	CancelBooking(bookingID string) error
	GetAllBookings() ([]domain.Booking, error)
	GetBookingsByRoomID(roomID string) ([]domain.Booking, error)
	GetBookingsByUserID(userID string) ([]domain.Booking, error)
	GetBookingsWithDetailsByRoomID(roomID string) ([]domain.BookingWithDetails, error)
	GetBookingsByDateRange(startDate, endDate int64) ([]domain.Booking, error)
	GetRoomScheduleByDate(roomID string, date int64) (*domain.RoomScheduleResponse, error)
}
