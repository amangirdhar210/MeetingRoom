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
	GetAll() ([]Room, error)
	GetByID(id int64) (*Room, error)
	UpdateAvailability(id int64, status string) error
	DeleteByID(id int64) error
	SearchWithFilters(minCapacity, maxCapacity int, floor *int, amenities []string) ([]Room, error)
}

type BookingRepository interface {
	Create(booking *Booking) error
	GetByID(id int64) (*Booking, error)
	GetAll() ([]Booking, error)
	GetByRoomAndTime(roomID int64, start, end time.Time) ([]Booking, error)
	GetByRoomID(roomID int64) ([]Booking, error)
	GetByUserID(userID int64) ([]Booking, error)
	Cancel(id int64) error
	GetByDateRange(startDate, endDate time.Time) ([]Booking, error)
	GetByRoomIDAndDate(roomID int64, date time.Time) ([]Booking, error)
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
	SearchRooms(minCapacity, maxCapacity int, floor *int, amenities []string, startTime, endTime *time.Time) ([]Room, error)
	CheckAvailability(roomID int64, startTime, endTime time.Time) (bool, []Booking, error)
	GetAvailableSlots(roomID int64, date time.Time, slotDuration int) ([]TimeSlot, error)
}

type BookingService interface {
	CreateBooking(booking *Booking) error
	GetBookingByID(bookingID int64) (*Booking, error)
	CancelBooking(bookingID int64) error
	GetAllBookings() ([]Booking, error)
	GetBookingsByRoomID(roomID int64) ([]Booking, error)
	GetBookingsByUserID(userID int64) ([]Booking, error)
	GetBookingsWithDetailsByRoomID(roomID int64) ([]BookingWithDetails, error)
	GetBookingsByDateRange(startDate, endDate time.Time) ([]Booking, error)
}

type TimeSlot struct {
	StartTime time.Time
	EndTime   time.Time
	Duration  int
}

type BookingWithDetails struct {
	Booking
	UserName   string
	UserEmail  string
	RoomName   string
	RoomNumber int
}
