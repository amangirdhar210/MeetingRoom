package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

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
	log.Println("CreateBooking handler invoked")
	log.Printf("Request body: %s", request.Body)
	log.Printf("Full authorizer context: %+v", request.RequestContext.Authorizer)

	var req dto.CreateBookingRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		log.Printf("Error unmarshaling request: %v", err)
		return shared.Response(400, dto.ErrorResponse{Error: "Invalid request body"})
	}

	log.Printf("Parsed request: RoomID=%s, StartTime=%s, EndTime=%s, Purpose=%s", req.RoomID, req.StartTime, req.EndTime, req.Purpose)

	var userID string
	if authContext, ok := request.RequestContext.Authorizer["lambda"].(map[string]any); ok {
		if uid, exists := authContext["userId"].(string); exists {
			userID = uid
		}
	}

	if userID == "" {
		log.Printf("User ID not found in authorizer context. Full context: %+v", request.RequestContext.Authorizer)
		return shared.Response(401, dto.ErrorResponse{Error: "Unauthorized"})
	}

	log.Printf("User ID from context: %s", userID)

	if req.StartTime == "" || req.EndTime == "" || req.RoomID == "" || req.Purpose == "" {
		log.Println("Missing required fields")
		return shared.Response(400, dto.ErrorResponse{Error: "All fields are required"})
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		log.Printf("Invalid start_time format: %v", err)
		return shared.Response(400, dto.ErrorResponse{Error: "Invalid start_time format"})
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		log.Printf("Invalid end_time format: %v", err)
		return shared.Response(400, dto.ErrorResponse{Error: "Invalid end_time format"})
	}

	log.Printf("Parsed times - Start: %d, End: %d", startTime.Unix(), endTime.Unix())

	booking := &domain.Booking{
		UserID:    userID,
		RoomID:    req.RoomID,
		StartTime: startTime.Unix(),
		EndTime:   endTime.Unix(),
		Purpose:   req.Purpose,
	}

	log.Printf("Creating booking: %+v", booking)

	if err := bookingService.CreateBooking(booking); err != nil {
		log.Printf("Error creating booking: %v", err)
		if err == domain.ErrRoomUnavailable {
			return shared.Response(409, dto.ErrorResponse{Error: "Room is already booked for this time"})
		}
		return shared.Response(500, dto.ErrorResponse{Error: err.Error()})
	}

	log.Printf("Booking created successfully with ID: %s", booking.ID)

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
