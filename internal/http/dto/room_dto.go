package dto

import "time"

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
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	RoomNumber  int      `json:"roomNumber"`
	Capacity    int      `json:"capacity"`
	Floor       int      `json:"floor"`
	Amenities   []string `json:"amenities"`
	Status      string   `json:"status"`
	Location    string   `json:"location"`
	Description string   `json:"description,omitempty"`
}

type RoomWithAvailabilityDTO struct {
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	RoomNumber       int        `json:"roomNumber"`
	Capacity         int        `json:"capacity"`
	Floor            int        `json:"floor"`
	Amenities        []string   `json:"amenities"`
	Status           string     `json:"status"`
	Location         string     `json:"location"`
	Description      string     `json:"description,omitempty"`
	IsAvailable      bool       `json:"isAvailable"`
	NextAvailableAt  *time.Time `json:"nextAvailableAt,omitempty"`
	CurrentBookingID *string    `json:"currentBookingId,omitempty"`
}

type RoomSearchFilters struct {
	MinCapacity int      `json:"minCapacity,omitempty"`
	MaxCapacity int      `json:"maxCapacity,omitempty"`
	Floor       *int     `json:"floor,omitempty"`
	Amenities   []string `json:"amenities,omitempty"`
	StartTime   string   `json:"startTime,omitempty"`
	EndTime     string   `json:"endTime,omitempty"`
	Available   *bool    `json:"available,omitempty"`
}

type AvailabilityCheckRequest struct {
	RoomID    string `json:"roomId" validate:"required"`
	StartTime string `json:"startTime" validate:"required"`
	EndTime   string `json:"endTime" validate:"required"`
}

type AvailabilityCheckResponse struct {
	Available        bool                    `json:"available"`
	RoomID           string                  `json:"roomId"`
	RoomName         string                  `json:"roomName"`
	RequestedStart   time.Time               `json:"requestedStart"`
	RequestedEnd     time.Time               `json:"requestedEnd"`
	ConflictingSlots []ConflictingBookingDTO `json:"conflictingSlots,omitempty"`
	SuggestedSlots   []TimeSlotDTO           `json:"suggestedSlots,omitempty"`
}

type TimeSlotDTO struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Duration  int       `json:"duration"`
}
