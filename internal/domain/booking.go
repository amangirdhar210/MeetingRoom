package domain

import "time"

type Booking struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	RoomID    string    `json:"room_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Purpose   string    `json:"purpose"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
