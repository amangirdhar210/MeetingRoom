package main

import (
	"context"

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

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	idToDelete := request.PathParameters["id"]
	if idToDelete == "" {
		return shared.Response(400, dto.ErrorResponse{Error: "User ID is required"})
	}

	var currentUserID string
	if authContext, ok := request.RequestContext.Authorizer["lambda"].(map[string]interface{}); ok {
		if uid, exists := authContext["userId"].(string); exists {
			currentUserID = uid
		}
	}

	if currentUserID == idToDelete {
		return shared.Response(403, dto.ErrorResponse{Error: "Cannot delete yourself"})
	}

	err := userService.DeleteUserByID(idToDelete)
	if err != nil {
		if err == domain.ErrNotFound {
			return shared.Response(404, dto.ErrorResponse{Error: "User not found"})
		}
		return shared.Response(500, dto.ErrorResponse{Error: "Failed to delete user"})
	}

	return shared.Response(200, dto.GenericResponse{Message: "User deleted successfully"})
}

func main() {
	lambda.Start(Handler)
}
