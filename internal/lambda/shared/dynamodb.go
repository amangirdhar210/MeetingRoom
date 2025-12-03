package shared

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var dynamoClient *dynamodb.Client
var tableName string
var jwtSecret string

func InitDynamoDB() (*dynamodb.Client, string, error) {
	if dynamoClient != nil {
		return dynamoClient, tableName, nil
	}

	tableName = os.Getenv("TABLE_NAME")
	if tableName == "" {
		tableName = "MeetingRoomSystem"
	}

	jwtSecret = os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key"
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, "", err
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)
	return dynamoClient, tableName, nil
}

func GetJWTSecret() string {
	return jwtSecret
}
