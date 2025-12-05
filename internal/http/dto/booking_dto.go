package dto

type CreateBookingRequest struct {
	UserID    string `json:"user_id"`
	RoomID    string `json:"room_id" validate:"required"`
	StartTime string `json:"start_time" validate:"required,datetime"`
	EndTime   string `json:"end_time" validate:"required,datetime"`
	Purpose   string `json:"purpose" validate:"required"`
}

type BookingDTO struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	RoomID    string `json:"room_id"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Purpose   string `json:"purpose"`
	Status    string `json:"status,omitempty"`
}

type DetailedBookingDTO struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	UserName   string `json:"userName"`
	UserEmail  string `json:"userEmail"`
	RoomID     string `json:"room_id"`
	RoomName   string `json:"roomName"`
	RoomNumber int    `json:"roomNumber"`
	StartTime  int64  `json:"start_time"`
	EndTime    int64  `json:"end_time"`
	Duration   int    `json:"duration"`
	Purpose    string `json:"purpose"`
	Status     string `json:"status"`
}

type ConflictingBookingDTO struct {
	BookingID string `json:"bookingId"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
	Purpose   string `json:"purpose,omitempty"`
}

type RoomScheduleResponse struct {
	RoomID     string               `json:"roomId"`
	RoomName   string               `json:"roomName"`
	Date       string               `json:"date"`
	Bookings   []DetailedBookingDTO `json:"bookings"`
	TotalSlots int                  `json:"totalSlots"`
}

type BookingDynamoDBItem struct {
	PK        string `dynamodbav:"PK"`
	SK        string `dynamodbav:"SK"`
	UserID    string `dynamodbav:"UserID"`
	RoomID    string `dynamodbav:"RoomID"`
	Date      int64  `dynamodbav:"Date"`
	ID        string `dynamodbav:"ID"`
	StartTime int64  `dynamodbav:"StartTime"`
	EndTime   int64  `dynamodbav:"EndTime"`
	Purpose   string `dynamodbav:"Purpose"`
	Status    string `dynamodbav:"Status"`
	CreatedAt int64  `dynamodbav:"CreatedAt"`
	UpdatedAt int64  `dynamodbav:"UpdatedAt"`
}
