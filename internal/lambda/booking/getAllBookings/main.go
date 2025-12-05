package main

import (
	"context"
	"log"

	dynamodbRepo "github.com/amangirdhar210/meeting-room/internal/adapters/repositories/dynamoDB"
	"github.com/amangirdhar210/meeting-room/internal/core/service"
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
	log.Println("GetAllBookings handler invoked")
	log.Printf("Full authorizer context: %+v", request.RequestContext.Authorizer)

	var userID, role string
	if authContext, ok := request.RequestContext.Authorizer["lambda"].(map[string]any); ok {
		if uid, exists := authContext["userId"].(string); exists {
			userID = uid
		}
		if r, exists := authContext["role"].(string); exists {
			role = r
		}
	}

	log.Printf("User: %s, Role: %s", userID, role)

	if role == "admin" {
		log.Println("Admin role detected, fetching all bookings")
		allBookings, getErr := bookingService.GetAllBookings()
		if getErr != nil {
			log.Printf("Error getting all bookings: %v", getErr)
			return shared.Response(500, map[string]string{"error": "Internal server error"})
		}
		log.Printf("Retrieved %d bookings", len(allBookings))
		return shared.Response(200, allBookings)
	}

	log.Println("User role detected, fetching user bookings")
	userBookings, getErr := bookingService.GetBookingsByUserID(userID)
	if getErr != nil {
		log.Printf("Error getting user bookings: %v", getErr)
		return shared.Response(500, map[string]string{"error": "Internal server error"})
	}
	log.Printf("Retrieved %d bookings for user", len(userBookings))
	return shared.Response(200, userBookings)
}

func main() {
	lambda.Start(handler)
}
