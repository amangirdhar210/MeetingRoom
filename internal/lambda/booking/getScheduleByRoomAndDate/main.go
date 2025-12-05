package main

import (
	"context"
	"log"
	"time"

	dynamodbRepo "github.com/amangirdhar210/meeting-room/internal/adapters/repositories/dynamoDB"
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
	log.Println("GetScheduleByRoomAndDate handler invoked")

	roomID := request.PathParameters["id"]
	if roomID == "" {
		log.Println("Room ID is missing")
		return shared.Response(400, dto.ErrorResponse{Error: "Invalid room id"})
	}

	dateStr := request.QueryStringParameters["date"]
	if dateStr == "" {
		log.Println("Date parameter is missing")
		return shared.Response(400, dto.ErrorResponse{Error: "Date parameter is required (format: YYYY-MM-DD)"})
	}

	targetDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		return shared.Response(400, dto.ErrorResponse{Error: "Invalid date format, use YYYY-MM-DD"})
	}

	schedule, err := bookingService.GetRoomScheduleByDate(roomID, targetDate.Unix())
	if err != nil {
		log.Printf("Error getting schedule by date: %v", err)
		return shared.Response(500, dto.ErrorResponse{Error: "Internal server error"})
	}

	log.Printf("Retrieved schedule for room %s on date %s", roomID, dateStr)
	return shared.Response(200, schedule)
}

func main() {
	lambda.Start(handler)
}
