package main

import (
	"context"
	"log"

	dynamodbRepo "github.com/amangirdhar210/meeting-room/internal/adapters/repositories/dynamoDB"
	"github.com/amangirdhar210/meeting-room/internal/core/domain"
	"github.com/amangirdhar210/meeting-room/internal/core/service"
	"github.com/amangirdhar210/meeting-room/internal/http/dto"
	"github.com/amangirdhar210/meeting-room/internal/lambda/shared"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var bookingService service.BookingService

func init() {
	dynamoClient, tableName, err := shared.InitDynamoDB()
	if err != nil {
		panic(err)
	}

	bookingRepo := dynamodbRepo.NewBookingRepositoryDynamoDB(dynamoClient, tableName)
	roomRepo := dynamodbRepo.NewRoomRepositoryDynamoDB(dynamoClient, tableName)
	userRepo := dynamodbRepo.NewUserRepositoryDynamoDB(dynamoClient, tableName)
	bookingService = service.NewBookingService(bookingRepo, roomRepo, userRepo)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("GetSchedule handler invoked")

	roomID := request.PathParameters["id"]
	if roomID == "" {
		log.Println("Room ID is missing")
		return shared.Response(400, dto.ErrorResponse{Error: "Invalid room id"})
	}

	detailedBookings, err := bookingService.GetBookingsWithDetailsByRoomID(roomID)
	if err != nil {
		log.Printf("Error getting schedule: %v", err)
		if err == domain.ErrNotFound {
			return shared.Response(200, []dto.DetailedBookingDTO{})
		}
		return shared.Response(500, dto.ErrorResponse{Error: "Internal server error"})
	}

	var response []dto.DetailedBookingDTO
	for _, booking := range detailedBookings {
		durationMinutes := int((booking.EndTime - booking.StartTime) / 60)
		response = append(response, dto.DetailedBookingDTO{
			ID:         booking.ID,
			UserID:     booking.UserID,
			UserName:   booking.UserName,
			UserEmail:  booking.UserEmail,
			RoomID:     booking.RoomID,
			RoomName:   booking.RoomName,
			RoomNumber: booking.RoomNumber,
			StartTime:  booking.StartTime,
			EndTime:    booking.EndTime,
			Duration:   durationMinutes,
			Purpose:    booking.Purpose,
			Status:     booking.Status,
		})
	}

	log.Printf("Retrieved %d bookings for room %s", len(response), roomID)
	return shared.Response(200, response)
}

func main() {
	lambda.Start(handler)
}
