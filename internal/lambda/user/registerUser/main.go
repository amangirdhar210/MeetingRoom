package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/amangirdhar210/meeting-room/internal/adapters/auth"
	dynamoRepo "github.com/amangirdhar210/meeting-room/internal/adapters/repositories/dynamoDB"
	"github.com/amangirdhar210/meeting-room/internal/core/domain"
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
	var registerReq dto.RegisterUserRequest
	if err := json.Unmarshal([]byte(request.Body), &registerReq); err != nil {
		return shared.Response(400, dto.ErrorResponse{Error: "Invalid request body"})
	}

	user := &domain.User{
		Name:     registerReq.Name,
		Email:    registerReq.Email,
		Password: registerReq.Password,
		Role:     registerReq.Role,
	}

	if err := userService.Register(user); err != nil {
		if err == domain.ErrConflict {
			return shared.Response(409, dto.ErrorResponse{Error: "User already exists"})
		}
		return shared.Response(500, dto.ErrorResponse{Error: "Failed to register user"})
	}

	return shared.Response(201, dto.UserDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	})
}

func main() {
	lambda.Start(handler)
}
