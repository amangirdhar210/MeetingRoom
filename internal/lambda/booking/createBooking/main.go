package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	dynamodbRepo "github.com/amangirdhar210/meeting-room/internal/adapters/repositories/dynamoDB"
	"github.com/amangirdhar210/meeting-room/internal/core/domain"
	"github.com/amangirdhar210/meeting-room/internal/core/service"
	"github.com/amangirdhar210/meeting-room/internal/http/dto"
	"github.com/amangirdhar210/meeting-room/internal/lambda/shared"
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
	var req dto.CreateBookingRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		log.Printf("Error unmarshaling request: %v", err)
		return shared.Response(400, dto.ErrorResponse{Error: "Invalid request body"})
	}

	userID := request.RequestContext.Authorizer["userId"].(string)
	if userID == "" {
		return shared.Response(401, dto.ErrorResponse{Error: "Unauthorized"})
	}

	if req.StartTime == "" || req.EndTime == "" || req.RoomID == "" || req.Purpose == "" {
		return shared.Response(400, dto.ErrorResponse{Error: "All fields are required"})
	}

	booking := &domain.Booking{
		UserID:    userID,
		RoomID:    req.RoomID,
		StartTime: 0,
		EndTime:   0,
		Purpose:   req.Purpose,
	}

	if err := bookingService.CreateBooking(booking); err != nil {
		log.Printf("Error creating booking: %v", err)
		return shared.Response(400, dto.ErrorResponse{Error: err.Error()})
	}

	return shared.Response(201, dto.BookingDTO{
		ID:        booking.ID,
		UserID:    booking.UserID,
		RoomID:    booking.RoomID,
		StartTime: booking.StartTime,
		EndTime:   booking.EndTime,
		Purpose:   booking.Purpose,
		Status:    booking.Status,
	})
}

func main() {
	lambda.Start(handler)
}
