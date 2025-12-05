package main

import (
	"context"
	"log"

	dynamodbRepo "github.com/amangirdhar210/meeting-room/internal/adapters/repositories/dynamoDB"
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
	log.Println("GetAllRooms handler invoked")

	rooms, err := roomService.GetAllRooms()
	if err != nil {
		log.Printf("Error getting rooms: %v", err)
		return shared.Response(500, dto.ErrorResponse{Error: "Internal server error"})
	}

	log.Printf("Retrieved %d rooms", len(rooms))
	return shared.Response(200, rooms)
}

func main() {
	lambda.Start(handler)
}
