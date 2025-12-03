package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/amangirdhar210/meeting-room/internal/adapters/auth"
	dynamoRepo "github.com/amangirdhar210/meeting-room/internal/adapters/repositories/dynamoDB"
	"github.com/amangirdhar210/meeting-room/internal/core/service"
	"github.com/amangirdhar210/meeting-room/internal/http/dto"
	"github.com/amangirdhar210/meeting-room/internal/lambda/shared"
)

var userService service.UserService

func init() {
	dynamoClient, tableName, err := shared.InitDynamoDB()
	if err != nil {
		panic(err)
	}

	userRepo := dynamoRepo.NewUserRepositoryDynamoDB(dynamoClient, tableName)

	hasher := auth.NewBcryptHasher()
	userService = service.NewUserService(userRepo, hasher)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := request.PathParameters["id"]
	if userID == "" {
		return shared.Response(400, dto.ErrorResponse{Error: "User ID is required"})
	}

	user, err := userService.GetUserByID(userID)
	if err != nil {
		return shared.Response(404, dto.ErrorResponse{Error: "User not found"})
	}

	return shared.Response(200, dto.UserDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func main() {
	lambda.Start(handler)
}
