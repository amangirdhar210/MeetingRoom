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
