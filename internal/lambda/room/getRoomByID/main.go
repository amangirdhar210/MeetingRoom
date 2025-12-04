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
	log.Println("GetRoomByID handler invoked")

	roomID := request.PathParameters["id"]
	if roomID == "" {
		log.Println("Room ID not provided in path parameters")
		return shared.Response(400, dto.ErrorResponse{Error: "Room ID is required"})
	}

	room, err := roomService.GetRoomByID(roomID)
	if err != nil {
		log.Printf("Error getting room by ID %s: %v", roomID, err)
		if err == domain.ErrNotFound {
			return shared.Response(404, dto.ErrorResponse{Error: "Room not found"})
		}
		if err == domain.ErrInvalidInput {
			return shared.Response(400, dto.ErrorResponse{Error: err.Error()})
		}
		return shared.Response(500, dto.ErrorResponse{Error: "Internal server error"})
	}

	log.Printf("Retrieved room with ID: %s", roomID)
	return shared.Response(200, room)
}

func main() {
	lambda.Start(handler)
}
