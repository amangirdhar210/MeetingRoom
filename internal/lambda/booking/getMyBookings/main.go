package main

import (
	"context"

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
	var userID string
	if authContext, ok := request.RequestContext.Authorizer["lambda"].(map[string]any); ok {
		if uid, exists := authContext["userId"].(string); exists {
			userID = uid
		}
	}

	if userID == "" {
		return shared.Response(401, map[string]string{"error": "Unauthorized"})
	}

	bookings, err := bookingService.GetBookingsByUserID(userID)
	if err != nil {
		return shared.Response(500, map[string]string{"error": "Internal server error"})
	}

	return shared.Response(200, bookings)
}

func main() {
	lambda.Start(handler)
}
