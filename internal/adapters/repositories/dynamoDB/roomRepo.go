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
		SK:          fmt.Sprintf("ROOM#%s", room.ID),
		LSI1:        room.Floor,
		LSI2:        room.Capacity,
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

	input := &dynamodb.GetItemInput{
		TableName: aws.String(repo.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "ROOM"},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("ROOM#%s", id)},
		},
	}

	result, err := repo.client.GetItem(ctx, input)
	if err != nil {
		log.Printf("Failed to get room with ID %s: %v", id, err)
		return nil, fmt.Errorf("failed to get room: %w", err)
	}

	if result.Item == nil {
		log.Printf("Room not found with ID: %s", id)
		return nil, domain.ErrNotFound
	}

	var roomItem dto.RoomDynamoDBItem
	err = attributevalue.UnmarshalMap(result.Item, &roomItem)
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

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(repo.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "ROOM"},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("ROOM#%s", id)},
		},
		ConditionExpression: aws.String("attribute_exists(PK) AND attribute_exists(SK)"),
	}

	_, err := repo.client.DeleteItem(ctx, input)
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
	ctx := context.Background()

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(repo.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "ROOM"},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("ROOM#%s", id)},
		},
		UpdateExpression: aws.String("SET #status = :status, UpdatedAt = :updatedAt"),
		ExpressionAttributeNames: map[string]string{
			"#status": "Status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status":    &types.AttributeValueMemberS{Value: status},
			":updatedAt": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", time.Now().Unix())},
		},
		ConditionExpression: aws.String("attribute_exists(PK) AND attribute_exists(SK)"),
	}

	_, err := repo.client.UpdateItem(ctx, input)
	if err != nil {
		log.Printf("Failed to update room availability: %v", err)
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to update room availability: %w", err)
	}

	log.Printf("Updated room availability: %s to %s", id, status)
	return nil
}

func (repo *RoomRepositoryDynamoDB) SearchWithFilters(minCapacity, maxCapacity int, floor *int) ([]domain.Room, error) {
	ctx := context.Background()

	if floor != nil {
		input := &dynamodb.QueryInput{
			TableName:              aws.String(repo.table),
			IndexName:              aws.String("LSI-1"),
			KeyConditionExpression: aws.String("PK = :pk AND LSI1 = :floor"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk":    &types.AttributeValueMemberS{Value: "ROOM"},
				":floor": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", *floor)},
			},
		}

		if minCapacity > 0 || maxCapacity > 0 {
			filterExprs := []string{}
			if minCapacity > 0 {
				filterExprs = append(filterExprs, "LSI2 >= :minCap")
				input.ExpressionAttributeValues[":minCap"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", minCapacity)}
			}
			if maxCapacity > 0 {
				filterExprs = append(filterExprs, "LSI2 <= :maxCap")
				input.ExpressionAttributeValues[":maxCap"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", maxCapacity)}
			}
			if len(filterExprs) > 0 {
				input.FilterExpression = aws.String(strings.Join(filterExprs, " AND "))
			}
		}

		result, err := repo.client.Query(ctx, input)
		if err != nil {
			log.Printf("Failed to search rooms: %v", err)
			return nil, fmt.Errorf("failed to search rooms: %w", err)
		}

		if len(result.Items) == 0 {
			return []domain.Room{}, nil
		}

		return repo.parseRoomItems(result.Items)
	}

	if minCapacity > 0 || maxCapacity > 0 {
		var keyConditionExpr string
		exprAttrValues := map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "ROOM"},
		}

		if minCapacity > 0 && maxCapacity > 0 {
			keyConditionExpr = "PK = :pk AND LSI2 BETWEEN :minCap AND :maxCap"
			exprAttrValues[":minCap"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", minCapacity)}
			exprAttrValues[":maxCap"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", maxCapacity)}
		} else if minCapacity > 0 {
			keyConditionExpr = "PK = :pk AND LSI2 >= :minCap"
			exprAttrValues[":minCap"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", minCapacity)}
		} else {
			keyConditionExpr = "PK = :pk AND LSI2 <= :maxCap"
			exprAttrValues[":maxCap"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", maxCapacity)}
		}

		input := &dynamodb.QueryInput{
			TableName:                 aws.String(repo.table),
			IndexName:                 aws.String("LSI-2"),
			KeyConditionExpression:    aws.String(keyConditionExpr),
			ExpressionAttributeValues: exprAttrValues,
		}

		result, err := repo.client.Query(ctx, input)
		if err != nil {
			log.Printf("Failed to search rooms: %v", err)
			return nil, fmt.Errorf("failed to search rooms: %w", err)
		}

		if len(result.Items) == 0 {
			return []domain.Room{}, nil
		}

		return repo.parseRoomItems(result.Items)
	}

	input := &dynamodb.QueryInput{
		TableName:              aws.String(repo.table),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "ROOM"},
		},
	}

	result, err := repo.client.Query(ctx, input)
	if err != nil {
		log.Printf("Failed to search rooms: %v", err)
		return nil, fmt.Errorf("failed to search rooms: %w", err)
	}

	if len(result.Items) == 0 {
		return []domain.Room{}, nil
	}

	return repo.parseRoomItems(result.Items)
}

func (repo *RoomRepositoryDynamoDB) parseRoomItems(items []map[string]types.AttributeValue) ([]domain.Room, error) {
	var rooms []domain.Room
	for _, item := range items {
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

	log.Printf("Search returned %d rooms", len(rooms))
	return rooms, nil
}
