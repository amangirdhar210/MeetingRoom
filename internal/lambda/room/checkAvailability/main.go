package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	dynamodbRepo "github.com/amangirdhar210/meeting-room/internal/adapters/repositories/dynamoDB"
	"github.com/amangirdhar210/meeting-room/internal/core/domain"
	"github.com/amangirdhar210/meeting-room/internal/core/service"
	"github.com/amangirdhar210/meeting-room/internal/http/dto"
	"github.com/amangirdhar210/meeting-room/internal/lambda/shared"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var roomService service.RoomService

func init() {
	dynamoClient, tableName, err := shared.InitDynamoDB()
	if err != nil {
		panic(err)
	}

	roomRepo := dynamodbRepo.NewRoomRepositoryDynamoDB(dynamoClient, tableName)
	roomService = service.NewRoomService(roomRepo)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("CheckAvailability handler invoked")

	var req dto.AvailabilityCheckRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		log.Printf("Error unmarshalling request body: %v", err)
		return shared.Response(400, dto.ErrorResponse{Error: "Invalid request body"})
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		log.Printf("Error parsing start time: %v", err)
		return shared.Response(400, dto.ErrorResponse{Error: "Invalid start_time format"})
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		log.Printf("Error parsing end time: %v", err)
		return shared.Response(400, dto.ErrorResponse{Error: "Invalid end_time format"})
	}

	isAvailable, conflictingBookings, err := roomService.CheckAvailability(req.RoomID, startTime.Unix(), endTime.Unix())
	if err != nil {
		log.Printf("Error checking availability: %v", err)
		if err == domain.ErrNotFound {
			return shared.Response(404, dto.ErrorResponse{Error: "Room not found"})
		}
		return shared.Response(500, dto.ErrorResponse{Error: "Internal server error"})
	}

	room, err := roomService.GetRoomByID(req.RoomID)
	if err != nil {
		log.Printf("Error getting room details: %v", err)
		return shared.Response(404, dto.ErrorResponse{Error: "Room not found"})
	}

	var conflictingSlots []dto.ConflictingBookingDTO
	for _, conflictBooking := range conflictingBookings {
		conflictingSlots = append(conflictingSlots, dto.ConflictingBookingDTO{
			BookingID: conflictBooking.ID,
			StartTime: conflictBooking.StartTime,
			EndTime:   conflictBooking.EndTime,
			Purpose:   conflictBooking.Purpose,
		})
	}

	response := dto.AvailabilityCheckResponse{
		Available:        isAvailable,
		RoomID:           req.RoomID,
		RoomName:         room.Name,
		RequestedStart:   startTime.Unix(),
		RequestedEnd:     endTime.Unix(),
		ConflictingSlots: conflictingSlots,
		SuggestedSlots:   []dto.TimeSlotDTO{},
	}

	log.Printf("Availability check completed - Available: %v", isAvailable)
	return shared.Response(200, response)
}

func main() {
	lambda.Start(handler)
}
