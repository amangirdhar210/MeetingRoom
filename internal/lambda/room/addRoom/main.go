package main

import (
	"context"
	"encoding/json"
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
	log.Println("AddRoom handler invoked")

	var req dto.AddRoomRequest

	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		log.Printf("Error unmarshalling request body: %v", err)
		return shared.Response(400, dto.ErrorResponse{Error: "Invalid request body"})
	}

	room := &domain.Room{
		Name:        req.Name,
		RoomNumber:  req.RoomNumber,
		Capacity:    req.Capacity,
		Floor:       req.Floor,
		Amenities:   req.Amenities,
		Status:      req.Status,
		Location:    req.Location,
		Description: req.Description,
	}

	err = roomService.AddRoom(room)
	if err != nil {
		log.Printf("Error adding room: %v", err)
		if err == domain.ErrInvalidInput {
			return shared.Response(400, dto.ErrorResponse{Error: err.Error()})
		}
		if err == domain.ErrConflict {
			return shared.Response(409, dto.ErrorResponse{Error: err.Error()})
		}
		return shared.Response(500, dto.ErrorResponse{Error: "Internal server error"})
	}

	log.Printf("Room created successfully with ID: %s", room.ID)
	return shared.Response(201, room)
}

func main() {
	lambda.Start(handler)
}
