package main

import (
	"context"
	"fmt"
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

func handler(ctx context.Context, request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := strings.TrimPrefix(request.AuthorizationToken, "Bearer ")

	if token == "" {
		fmt.Println("Unauthorized: No token provided")
		return denyPolicy(request.MethodArn), nil
	}

	claims, err := jwtGenerator.ValidateToken(token)
	if err != nil {
		fmt.Printf("Unauthorized: Invalid token - %v\n", err)
		return denyPolicy(request.MethodArn), nil
	}

	if claims.UserID == "" {
		fmt.Println("Unauthorized: Invalid user ID in token")
		return denyPolicy(request.MethodArn), nil
	}

	if claims.Role != "admin" {
		fmt.Printf("Unauthorized: User %s is not admin (role: %s)\n", claims.UserID, claims.Role)
		return denyPolicy(request.MethodArn), nil
	}

	return allowPolicy(request.MethodArn, claims.UserID, claims.Role), nil
}

func allowPolicy(resource, userID, role string) events.APIGatewayCustomAuthorizerResponse {
	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: userID,
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   "Allow",
					Resource: []string{resource},
				},
			},
		},
		Context: map[string]any{
			"userId": userID,
			"role":   role,
		},
	}
}

func denyPolicy(resource string) events.APIGatewayCustomAuthorizerResponse {
	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: "unauthorized",
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   "Deny",
					Resource: []string{resource},
				},
			},
		},
	}
}

func main() {
	lambda.Start(handler)
}
