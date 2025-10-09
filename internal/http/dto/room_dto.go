package dto

type AddRoomRequest struct {
	Name        string `json:"name" validate:"required"`
	Capacity    int    `json:"capacity" validate:"required,min=1"`
	Location    string `json:"location" validate:"required"`
	Description string `json:"description,omitempty"`
}

type RoomDTO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Capacity    int    `json:"capacity"`
	Location    string `json:"location"`
	Description string `json:"description,omitempty"`
	IsAvailable bool   `json:"is_available"`
}
