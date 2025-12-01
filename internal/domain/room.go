package domain

import "time"

type Room struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	RoomNumber  int       `json:"roomNumber"`
	Capacity    int       `json:"capacity"`
	Floor       int       `json:"floor"`
	Amenities   []string  `json:"amenities"`
	Status      string    `json:"status"`
	Location    string    `json:"location"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
