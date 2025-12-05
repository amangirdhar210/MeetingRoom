package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/amangirdhar210/meeting-room/internal/adapters/auth"
)

var jwtGenerator *auth.JWTGenerator

func init() {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key"
	}
	jwtGenerator = auth.NewJWTGenerator(jwtSecret, 0)
}

func handler(ctx context.Context, request events.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	authHeader := request.Headers["authorization"]
	if authHeader == "" {
		authHeader = request.Headers["Authorization"]
	}

	if authHeader == "" {
		log.Printf("Unauthorized: No Authorization header")
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
		}, nil
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" || token == authHeader {
		log.Printf("Unauthorized: Invalid Bearer token format")
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
		}, nil
	}

	claims, err := jwtGenerator.ValidateToken(token)
	if err != nil {
		log.Printf("Unauthorized: Invalid token - %v\n", err)
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
		}, nil
	}

	if claims.UserID == "" {
		log.Println("Unauthorized: Invalid user ID in token")
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
		}, nil
	}

	if claims.Role != "admin" {
		log.Printf("Unauthorized: User %s is not admin (role: %s)\n", claims.UserID, claims.Role)
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
		}, nil
	}

	log.Printf("Authorized: Admin user %s\n", claims.UserID)
	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: true,
		Context: map[string]any{
			"userId": claims.UserID,
			"role":   claims.Role,
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
