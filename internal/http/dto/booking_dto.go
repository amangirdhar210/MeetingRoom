package dto

import "time"

type CreateBookingRequest struct {
	UserID    int64  `json:"user_id"`
	RoomID    int64  `json:"room_id" validate:"required"`
	StartTime string `json:"start_time" validate:"required,datetime"`
	EndTime   string `json:"end_time" validate:"required,datetime"`
	Purpose   string `json:"purpose" validate:"required"`
}

type BookingDTO struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	RoomID    int64     `json:"room_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Purpose   string    `json:"purpose"`
	Status    string    `json:"status,omitempty"`
}

type DetailedBookingDTO struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	UserName   string    `json:"userName"`
	UserEmail  string    `json:"userEmail"`
	RoomID     int64     `json:"room_id"`
	RoomName   string    `json:"roomName"`
	RoomNumber int       `json:"roomNumber"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Duration   int       `json:"duration"`
	Purpose    string    `json:"purpose"`
	Status     string    `json:"status"`
}

type ConflictingBookingDTO struct {
	BookingID int64     `json:"bookingId"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Purpose   string    `json:"purpose,omitempty"`
}

type RoomScheduleResponse struct {
	RoomID     int64                `json:"roomId"`
	RoomName   string               `json:"roomName"`
	Date       string               `json:"date"`
	Bookings   []DetailedBookingDTO `json:"bookings"`
	TotalSlots int                  `json:"totalSlots"`
}
