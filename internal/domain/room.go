package domain

import "time"

type Room struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Capacity    int       `json:"capacity"`
	Location    string    `json:"location"`
	Description string    `json:"description,omitempty"`
	IsAvailable bool      `json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
