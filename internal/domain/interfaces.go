package domain

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	GetByID(id string) (*User, error)
	GetAll() ([]User, error)
	DeleteByID(id string) error
}

type RoomRepository interface {
	Create(room *Room) error
	GetAll() ([]Room, error)
	GetByID(id string) (*Room, error)
	UpdateAvailability(id string, status string) error
	DeleteByID(id string) error
	SearchWithFilters(minCapacity, maxCapacity int, floor *int, amenities []string) ([]Room, error)
}

type BookingRepository interface {
	Create(booking *Booking) error
	GetByID(id string) (*Booking, error)
	GetAll() ([]Booking, error)
	GetByRoomAndTime(roomID string, start, end int64) ([]Booking, error)
	GetByRoomID(roomID string) ([]Booking, error)
	GetByUserID(userID string) ([]Booking, error)
	Cancel(id string) error
	GetByDateRange(startDate, endDate int64) ([]Booking, error)
	GetByRoomIDAndDate(roomID string, date int64) ([]Booking, error)
}

type UserService interface {
	Register(user *User) error
	GetAllUsers() ([]User, error)
	GetUserByID(id string) (*User, error)
	DeleteUserByID(id string) error
}

type AuthService interface {
	HashPassword(password string) (string, error)

	VerifyPassword(hashed, plain string) bool

	AuthenticateUser(email, password string) (token string, user *User, err error)
}

type RoomService interface {
	AddRoom(room *Room) error
	GetAllRooms() ([]Room, error)
	GetRoomByID(id string) (*Room, error)
	DeleteRoomByID(id string) error
	SearchRooms(minCapacity, maxCapacity int, floor *int, amenities []string, startTime, endTime *int64) ([]Room, error)
	CheckAvailability(roomID string, startTime, endTime int64) (bool, []Booking, error)
	GetAvailableSlots(roomID string, date int64, slotDuration int) ([]TimeSlot, error)
}

type BookingService interface {
	CreateBooking(booking *Booking) error
	GetBookingByID(bookingID string) (*Booking, error)
	CancelBooking(bookingID string) error
	GetAllBookings() ([]Booking, error)
	GetBookingsByRoomID(roomID string) ([]Booking, error)
	GetBookingsByUserID(userID string) ([]Booking, error)
	GetBookingsWithDetailsByRoomID(roomID string) ([]BookingWithDetails, error)
	GetBookingsByDateRange(startDate, endDate int64) ([]Booking, error)
	GetRoomScheduleByDate(roomID string, date int64) (*RoomScheduleResponse, error)
}

type TimeSlot struct {
	StartTime int64
	EndTime   int64
	Duration  int
}

type BookingWithDetails struct {
	Booking
	UserName   string
	UserEmail  string
	RoomName   string
	RoomNumber int
}

type ScheduleSlot struct {
	StartTime string  `json:"startTime"`
	EndTime   string  `json:"endTime"`
	IsBooked  bool    `json:"isBooked"`
	BookingID *string `json:"bookingId,omitempty"`
	UserName  string  `json:"userName,omitempty"`
	Purpose   string  `json:"purpose,omitempty"`
}

type RoomScheduleResponse struct {
	RoomID     string         `json:"roomId"`
	RoomName   string         `json:"roomName"`
	RoomNumber int            `json:"roomNumber"`
	Date       string         `json:"date"`
	Bookings   []ScheduleSlot `json:"bookings"`
}
