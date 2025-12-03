package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/amangirdhar210/meeting-room/internal/lambda/shared"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return shared.Response(200, map[string]string{
		"status":  "healthy",
		"service": "MeetingRoom API",
	})
}

func main() {
	lambda.Start(handler)
}
