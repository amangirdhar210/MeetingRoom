package domain

type Booking struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	RoomID    string `json:"room_id"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Purpose   string `json:"purpose"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
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
