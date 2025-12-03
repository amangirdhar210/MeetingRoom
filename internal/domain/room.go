package domain

type Room struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	RoomNumber  int      `json:"roomNumber"`
	Capacity    int      `json:"capacity"`
	Floor       int      `json:"floor"`
	Amenities   []string `json:"amenities"`
	Status      string   `json:"status"`
	Location    string   `json:"location"`
	Description string   `json:"description,omitempty"`
	CreatedAt   int64    `json:"created_at"`
	UpdatedAt   int64    `json:"updated_at"`
}
