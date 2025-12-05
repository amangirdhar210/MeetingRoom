package dynamodb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/amangirdhar210/meeting-room/internal/core/domain"
	"github.com/amangirdhar210/meeting-room/internal/core/ports"
	"github.com/amangirdhar210/meeting-room/internal/http/dto"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type RoomRepositoryDynamoDB struct {
	client *dynamodb.Client
	table  string
}

func NewRoomRepositoryDynamoDB(client *dynamodb.Client, tableName string) ports.RoomRepository {
	return &RoomRepositoryDynamoDB{
		client: client,
		table:  tableName,
	}
}

func (repo *RoomRepositoryDynamoDB) Create(room *domain.Room) error {
	ctx := context.Background()

	if room.ID == "" {
		room.ID = uuid.New().String()
	}

	now := time.Now().Unix()
	room.CreatedAt = now
	room.UpdatedAt = now

	if room.Status == "" {
		room.Status = "available"
	}

	item := dto.RoomDynamoDBItem{
		PK:          "ROOM",
		SK:          fmt.Sprintf("ROOM#%d", room.Capacity),
		ID:          room.ID,
		Name:        room.Name,
		RoomNumber:  room.RoomNumber,
		Capacity:    room.Capacity,
		Floor:       room.Floor,
		Amenities:   room.Amenities,
		Status:      room.Status,
		Location:    room.Location,
		Description: room.Description,
		CreatedAt:   room.CreatedAt,
		UpdatedAt:   room.UpdatedAt,
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		log.Printf("Failed to marshal room: %v", err)
		return fmt.Errorf("failed to marshal room: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName:           aws.String(repo.table),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
	}

	_, err = repo.client.PutItem(ctx, input)
	if err != nil {
		log.Printf("Failed to create room: %v", err)
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return domain.ErrConflict
		}
		return fmt.Errorf("failed to create room: %w", err)
	}

	log.Printf("Room created successfully with ID: %s", room.ID)
	return nil
}

func (repo *RoomRepositoryDynamoDB) GetAll() ([]domain.Room, error) {
	ctx := context.Background()

	input := &dynamodb.QueryInput{
		TableName:              aws.String(repo.table),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "ROOM"},
		},
	}

	result, err := repo.client.Query(ctx, input)
	if err != nil {
		log.Printf("Failed to query rooms: %v", err)
		return nil, fmt.Errorf("failed to query rooms: %w", err)
	}

	if len(result.Items) == 0 {
		return []domain.Room{}, nil
	}

	var rooms []domain.Room
	for _, item := range result.Items {
		var roomItem dto.RoomDynamoDBItem
		err := attributevalue.UnmarshalMap(item, &roomItem)
		if err != nil {
			log.Printf("Failed to unmarshal room item: %v", err)
			continue
		}

		room := domain.Room{
			ID:          roomItem.ID,
			Name:        roomItem.Name,
			RoomNumber:  roomItem.RoomNumber,
			Capacity:    roomItem.Capacity,
			Floor:       roomItem.Floor,
			Amenities:   roomItem.Amenities,
			Status:      roomItem.Status,
			Location:    roomItem.Location,
			Description: roomItem.Description,
			CreatedAt:   roomItem.CreatedAt,
			UpdatedAt:   roomItem.UpdatedAt,
		}
		rooms = append(rooms, room)
	}

	log.Printf("Retrieved %d rooms", len(rooms))
	return rooms, nil
}

func (repo *RoomRepositoryDynamoDB) GetByID(id string) (*domain.Room, error) {
	ctx := context.Background()

	input := &dynamodb.ScanInput{
		TableName:        aws.String(repo.table),
		FilterExpression: aws.String("PK = :pk AND ID = :id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "ROOM"},
			":id": &types.AttributeValueMemberS{Value: id},
		},
	}

	result, err := repo.client.Scan(ctx, input)
	if err != nil {
		log.Printf("Failed to scan for room with ID %s: %v", id, err)
		return nil, fmt.Errorf("failed to get room: %w", err)
	}

	if len(result.Items) == 0 {
		log.Printf("Room not found with ID: %s", id)
		return nil, domain.ErrNotFound
	}

	var roomItem dto.RoomDynamoDBItem
	err = attributevalue.UnmarshalMap(result.Items[0], &roomItem)
	if err != nil {
		log.Printf("Failed to unmarshal room: %v", err)
		return nil, fmt.Errorf("failed to unmarshal room: %w", err)
	}

	room := &domain.Room{
		ID:          roomItem.ID,
		Name:        roomItem.Name,
		RoomNumber:  roomItem.RoomNumber,
		Capacity:    roomItem.Capacity,
		Floor:       roomItem.Floor,
		Amenities:   roomItem.Amenities,
		Status:      roomItem.Status,
		Location:    roomItem.Location,
		Description: roomItem.Description,
		CreatedAt:   roomItem.CreatedAt,
		UpdatedAt:   roomItem.UpdatedAt,
	}

	log.Printf("Retrieved room with ID: %s", id)
	return room, nil
}

func (repo *RoomRepositoryDynamoDB) DeleteByID(id string) error {
	ctx := context.Background()

	room, err := repo.GetByID(id)
	if err != nil {
		return err
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(repo.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "ROOM"},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("ROOM#%d", room.Capacity)},
		},
		ConditionExpression: aws.String("attribute_exists(PK) AND attribute_exists(SK)"),
	}

	_, err = repo.client.DeleteItem(ctx, input)
	if err != nil {
		log.Printf("Failed to delete room with ID %s: %v", id, err)
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to delete room: %w", err)
	}

	log.Printf("Deleted room with ID: %s", id)
	return nil
}

func (repo *RoomRepositoryDynamoDB) UpdateAvailability(id string, status string) error {
	return nil
}

func (repo *RoomRepositoryDynamoDB) SearchWithFilters(minCapacity, maxCapacity int, floor *int, amenities []string) ([]domain.Room, error) {
	return nil, nil
}
