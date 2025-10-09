package domain

import "time"

//
// ────────────────────────────────
//   REPOSITORY INTERFACES
//   (Match repo usage in service implementations)
// ────────────────────────────────
//

// UserRepository defines persistence operations for users.
type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	GetByID(id int64) (*User, error)
	GetAll() ([]User, error)
}

// RoomRepository defines persistence operations for rooms.
type RoomRepository interface {
	Create(room *Room) error
	GetByID(id int64) (*Room, error)
	GetAll() ([]Room, error)
	UpdateAvailability(id int64, available bool) error
}

// BookingRepository defines persistence operations for bookings.
type BookingRepository interface {
	Create(booking *Booking) error
	GetByID(id int64) (*Booking, error)
	GetAll() ([]Booking, error)
	// Returns bookings for a room that overlap (or fall in) the provided time range.
	GetByRoomAndTime(roomID int64, start, end time.Time) ([]Booking, error)
	Cancel(id int64) error
}

//
// ────────────────────────────────
//   SERVICE INTERFACES
//   (Match service implementations present in internal/service)
// ────────────────────────────────
//

// UserService handles user management (registration, listing, etc.)
type UserService interface {
	// Register expects a populated domain.User (hashed password handled by AuthService in implementations)
	Register(user *User) error
	GetAllUsers() ([]User, error)
}

// AuthService handles authentication concerns:
// - password hashing/verification
// - authenticating credentials (returns token + user in our implementation)
type AuthService interface {
	// HashPassword hashes a plaintext password.
	HashPassword(password string) (string, error)

	// VerifyPassword compares stored hash with plaintext password.
	VerifyPassword(hashed, plain string) bool

	// AuthenticateUser verifies credentials and returns a JWT token and the user.
	AuthenticateUser(email, password string) (token string, user *User, err error)
}

// RoomService handles room-related business logic.
type RoomService interface {
	AddRoom(room *Room) error
	GetAllRooms() ([]Room, error)
	GetRoomByID(id int64) (*Room, error)
}

// BookingService handles booking-related business logic.
type BookingService interface {
	CreateBooking(booking *Booking) error
	CancelBooking(id int64) error
	GetAllBookings() ([]Booking, error)
}
