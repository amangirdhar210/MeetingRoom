package domain

import "time"

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	GetByID(id int64) (*User, error)
	GetAll() ([]User, error)
	DeleteByID(id int64) error
}

type RoomRepository interface {
	Create(room *Room) error
	GetByID(id int64) (*Room, error)
	GetAll() ([]Room, error)
	UpdateAvailability(id int64, available bool) error
	DeleteByID(id int64) error
}

type BookingRepository interface {
	Create(booking *Booking) error
	GetByID(id int64) (*Booking, error)
	GetAll() ([]Booking, error)
	GetByRoomAndTime(roomID int64, start, end time.Time) ([]Booking, error)
	Cancel(id int64) error
}

type UserService interface {
	Register(user *User) error
	GetAllUsers() ([]User, error)
	GetUserByID(id int64) (*User, error)
	DeleteUserByID(id int64) error
}

type AuthService interface {
	HashPassword(password string) (string, error)

	VerifyPassword(hashed, plain string) bool

	AuthenticateUser(email, password string) (token string, user *User, err error)
}

type RoomService interface {
	AddRoom(room *Room) error
	GetAllRooms() ([]Room, error)
	GetRoomByID(id int64) (*Room, error)
	DeleteRoomByID(id int64) error
}

type BookingService interface {
	CreateBooking(booking *Booking) error
	CancelBooking(id int64) error
	GetAllBookings() ([]Booking, error)
}
