package dto

import "time"

type CreateBookingRequest struct {
	UserID    int64  `json:"user_id"` // handler injects userID from JWT
	RoomID    int64  `json:"room_id" validate:"required"`
	StartTime string `json:"start_time" validate:"required,datetime"` // RFC3339
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
