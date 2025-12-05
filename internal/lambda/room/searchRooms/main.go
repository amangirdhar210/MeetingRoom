package main

import (
	"context"
	"log"
	"strconv"

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
	log.Println("SearchRooms handler invoked")

	queryParams := request.QueryStringParameters

	minCapacity := 0
	if minCapStr, ok := queryParams["minCapacity"]; ok && minCapStr != "" {
		if val, err := strconv.Atoi(minCapStr); err == nil {
			minCapacity = val
		}
	}

	maxCapacity := 0
	if maxCapStr, ok := queryParams["maxCapacity"]; ok && maxCapStr != "" {
		if val, err := strconv.Atoi(maxCapStr); err == nil {
			maxCapacity = val
		}
	}

	var floor *int
	if floorStr, ok := queryParams["floor"]; ok && floorStr != "" {
		if val, err := strconv.Atoi(floorStr); err == nil {
			floor = &val
		}
	}

	rooms, err := roomService.SearchRooms(minCapacity, maxCapacity, floor, nil, nil)
	if err != nil {
		log.Printf("Error searching rooms: %v", err)
		return shared.Response(500, dto.ErrorResponse{Error: "Internal server error"})
	}

	log.Printf("Found %d rooms matching search criteria", len(rooms))
	return shared.Response(200, rooms)
}

func main() {
	lambda.Start(handler)
}
