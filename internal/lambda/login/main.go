package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/amangirdhar210/meeting-room/internal/adapters/auth"
	dynamoRepo "github.com/amangirdhar210/meeting-room/internal/adapters/repositories/dynamoDB"
	"github.com/amangirdhar210/meeting-room/internal/core/ports"
	"github.com/amangirdhar210/meeting-room/internal/core/service"
	"github.com/amangirdhar210/meeting-room/internal/http/dto"
	"github.com/amangirdhar210/meeting-room/internal/lambda/shared"
)

var authService service.AuthService
var jwtGen ports.TokenGenerator
var hasher ports.PasswordHasher

func init() {
	dynamoClient, tableName, err := shared.InitDynamoDB()
	if err != nil {
		panic(err)
	}

	userRepo := dynamoRepo.NewUserRepositoryDynamoDB(dynamoClient, tableName)

	jwtGen = auth.NewJWTGenerator(shared.GetJWTSecret(), 24*time.Hour)
	hasher = auth.NewBcryptHasher()
	authService = service.NewAuthService(userRepo, jwtGen, hasher)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var loginReq dto.LoginUserRequest
	if err := json.Unmarshal([]byte(request.Body), &loginReq); err != nil {
		return shared.Response(400, dto.ErrorResponse{Error: "Invalid request body"})
	}

	log.Printf("%+v", loginReq)

	token, user, err := authService.Login(loginReq.Email, loginReq.Password)
	if err != nil {
		return shared.Response(401, dto.ErrorResponse{Error: "Invalid credentials"})
	}

	return shared.Response(200, dto.LoginUserResponse{
		Token: token,
		User: dto.UserDTO{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	})
}

func main() {
	lambda.Start(handler)
}
