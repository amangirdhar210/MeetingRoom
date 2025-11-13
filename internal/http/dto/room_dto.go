package dto

type AddRoomRequest struct {
	Name        string   `json:"name" validate:"required"`
	RoomNumber  int      `json:"roomNumber" validate:"required"`
	Capacity    int      `json:"capacity" validate:"required,min=1"`
	Floor       int      `json:"floor" validate:"required"`
	Amenities   []string `json:"amenities"`
	Status      string   `json:"status"`
	Location    string   `json:"location" validate:"required"`
	Description string   `json:"description,omitempty"`
}

type RoomDTO struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	RoomNumber  int      `json:"roomNumber"`
	Capacity    int      `json:"capacity"`
	Floor       int      `json:"floor"`
	Amenities   []string `json:"amenities"`
	Status      string   `json:"status"`
	Location    string   `json:"location"`
	Description string   `json:"description,omitempty"`
}
