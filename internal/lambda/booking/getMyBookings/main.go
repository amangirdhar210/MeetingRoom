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

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{}, nil
}

func main() {
	lambda.Start(handler)
}
