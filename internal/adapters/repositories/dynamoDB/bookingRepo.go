package dynamodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"

	"github.com/amangirdhar210/meeting-room/internal/core/domain"
	"github.com/amangirdhar210/meeting-room/internal/core/ports"
	"github.com/amangirdhar210/meeting-room/internal/http/dto"
)

type BookingRepositoryDynamoDB struct {
	client *dynamodb.Client
	table  string
}

func NewBookingRepositoryDynamoDB(client *dynamodb.Client, tableName string) ports.BookingRepository {
	return &BookingRepositoryDynamoDB{
		client: client,
		table:  tableName,
	}
}

func (repo *BookingRepositoryDynamoDB) Create(booking *domain.Booking) error {
	ctx := context.Background()

	if booking.ID == "" {
		booking.ID = uuid.New().String()
	}

	now := time.Now().Unix()
	booking.CreatedAt = now
	booking.UpdatedAt = now

	if booking.Status == "" {
		booking.Status = "confirmed"
	}

	startOfDay := (booking.StartTime / 86400) * 86400

	item := dto.BookingDynamoDBItem{
		PK:        "BOOKING",
		SK:        fmt.Sprintf("BOOKING#%s", booking.ID),
		UserID:    booking.UserID,
		RoomID:    booking.RoomID,
		Date:      startOfDay,
		ID:        booking.ID,
		StartTime: booking.StartTime,
		EndTime:   booking.EndTime,
		Purpose:   booking.Purpose,
		Status:    booking.Status,
		CreatedAt: booking.CreatedAt,
		UpdatedAt: booking.UpdatedAt,
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		log.Printf("Failed to marshal booking: %v", err)
		return fmt.Errorf("failed to marshal booking: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(repo.table),
		Item:      av,
	}

	_, err = repo.client.PutItem(ctx, input)
	if err != nil {
		log.Printf("Failed to create booking: %v", err)
		return fmt.Errorf("failed to create booking: %w", err)
	}

	log.Printf("Booking created successfully with ID: %s", booking.ID)
	return nil
}

func (repo *BookingRepositoryDynamoDB) GetByID(id string) (*domain.Booking, error) {
	ctx := context.Background()

	input := &dynamodb.GetItemInput{
		TableName: aws.String(repo.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "BOOKING"},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("BOOKING#%s", id)},
		},
	}

	result, err := repo.client.GetItem(ctx, input)
	if err != nil {
		log.Printf("Failed to get booking: %v", err)
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	if result.Item == nil {
		return nil, domain.ErrNotFound
	}

	var item dto.BookingDynamoDBItem
	err = attributevalue.UnmarshalMap(result.Item, &item)
	if err != nil {
		log.Printf("Failed to unmarshal booking: %v", err)
		return nil, fmt.Errorf("failed to unmarshal booking: %w", err)
	}

	booking := &domain.Booking{
		ID:        item.ID,
		UserID:    item.UserID,
		RoomID:    item.RoomID,
		StartTime: item.StartTime,
		EndTime:   item.EndTime,
		Purpose:   item.Purpose,
		Status:    item.Status,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}

	return booking, nil
}

func (repo *BookingRepositoryDynamoDB) GetAll() ([]domain.Booking, error) {
	ctx := context.Background()

	input := &dynamodb.QueryInput{
		TableName:              aws.String(repo.table),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "BOOKING"},
		},
	}

	result, err := repo.client.Query(ctx, input)
	if err != nil {
		log.Printf("Failed to get all bookings: %v", err)
		return nil, fmt.Errorf("failed to get all bookings: %w", err)
	}

	if len(result.Items) == 0 {
		return []domain.Booking{}, nil
	}

	var items []dto.BookingDynamoDBItem
	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		log.Printf("Failed to unmarshal bookings: %v", err)
		return nil, fmt.Errorf("failed to unmarshal bookings: %w", err)
	}

	bookings := make([]domain.Booking, len(items))
	for i, item := range items {
		bookings[i] = domain.Booking{
			ID:        item.ID,
			UserID:    item.UserID,
			RoomID:    item.RoomID,
			StartTime: item.StartTime,
			EndTime:   item.EndTime,
			Purpose:   item.Purpose,
			Status:    item.Status,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}

	return bookings, nil
}

func (repo *BookingRepositoryDynamoDB) GetByRoomAndTime(roomID string, start, end int64) ([]domain.Booking, error) {
	ctx := context.Background()

	input := &dynamodb.QueryInput{
		TableName:              aws.String(repo.table),
		IndexName:              aws.String("LSI-5"),
		KeyConditionExpression: aws.String("PK = :pk AND RoomID = :roomId"),
		FilterExpression:       aws.String("EndTime > :start AND StartTime < :end"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: "BOOKING"},
			":roomId": &types.AttributeValueMemberS{Value: roomID},
			":start":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", start)},
			":end":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", end)},
		},
	}

	result, err := repo.client.Query(ctx, input)
	if err != nil {
		log.Printf("Failed to get bookings by room and time: %v", err)
		return nil, fmt.Errorf("failed to get bookings by room and time: %w", err)
	}

	if len(result.Items) == 0 {
		return []domain.Booking{}, nil
	}

	var items []dto.BookingDynamoDBItem
	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		log.Printf("Failed to unmarshal bookings: %v", err)
		return nil, fmt.Errorf("failed to unmarshal bookings: %w", err)
	}

	bookings := make([]domain.Booking, len(items))
	for i, item := range items {
		bookings[i] = domain.Booking{
			ID:        item.ID,
			UserID:    item.UserID,
			RoomID:    item.RoomID,
			StartTime: item.StartTime,
			EndTime:   item.EndTime,
			Purpose:   item.Purpose,
			Status:    item.Status,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}

	return bookings, nil
}

func (repo *BookingRepositoryDynamoDB) GetByRoomID(roomID string) ([]domain.Booking, error) {
	ctx := context.Background()

	input := &dynamodb.QueryInput{
		TableName:              aws.String(repo.table),
		IndexName:              aws.String("LSI-5"),
		KeyConditionExpression: aws.String("PK = :pk AND RoomID = :roomId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: "BOOKING"},
			":roomId": &types.AttributeValueMemberS{Value: roomID},
		},
	}

	result, err := repo.client.Query(ctx, input)
	if err != nil {
		log.Printf("Failed to get bookings by room: %v", err)
		return nil, fmt.Errorf("failed to get bookings by room: %w", err)
	}

	if len(result.Items) == 0 {
		return []domain.Booking{}, nil
	}

	var items []dto.BookingDynamoDBItem
	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		log.Printf("Failed to unmarshal bookings: %v", err)
		return nil, fmt.Errorf("failed to unmarshal bookings: %w", err)
	}

	bookings := make([]domain.Booking, len(items))
	for i, item := range items {
		bookings[i] = domain.Booking{
			ID:        item.ID,
			UserID:    item.UserID,
			RoomID:    item.RoomID,
			StartTime: item.StartTime,
			EndTime:   item.EndTime,
			Purpose:   item.Purpose,
			Status:    item.Status,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}

	return bookings, nil
}

func (repo *BookingRepositoryDynamoDB) GetByUserID(userID string) ([]domain.Booking, error) {
	ctx := context.Background()

	input := &dynamodb.QueryInput{
		TableName:              aws.String(repo.table),
		IndexName:              aws.String("LSI-3"),
		KeyConditionExpression: aws.String("PK = :pk AND UserID = :userId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: "BOOKING"},
			":userId": &types.AttributeValueMemberS{Value: userID},
		},
	}

	result, err := repo.client.Query(ctx, input)
	if err != nil {
		log.Printf("Failed to get bookings by user: %v", err)
		return nil, fmt.Errorf("failed to get bookings by user: %w", err)
	}

	if len(result.Items) == 0 {
		return []domain.Booking{}, nil
	}

	var items []dto.BookingDynamoDBItem
	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		log.Printf("Failed to unmarshal bookings: %v", err)
		return nil, fmt.Errorf("failed to unmarshal bookings: %w", err)
	}

	bookings := make([]domain.Booking, len(items))
	for i, item := range items {
		bookings[i] = domain.Booking{
			ID:        item.ID,
			UserID:    item.UserID,
			RoomID:    item.RoomID,
			StartTime: item.StartTime,
			EndTime:   item.EndTime,
			Purpose:   item.Purpose,
			Status:    item.Status,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}

	return bookings, nil
}

func (repo *BookingRepositoryDynamoDB) Cancel(id string) error {
	ctx := context.Background()

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(repo.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "BOOKING"},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("BOOKING#%s", id)},
		},
		ConditionExpression: aws.String("attribute_exists(PK) AND attribute_exists(SK)"),
	}

	_, err := repo.client.DeleteItem(ctx, input)
	if err != nil {
		log.Printf("Failed to delete booking: %v", err)
		return fmt.Errorf("failed to delete booking: %w", err)
	}

	log.Printf("Booking deleted successfully: %s", id)
	return nil
}

func (repo *BookingRepositoryDynamoDB) GetByDateRange(startDate, endDate int64) ([]domain.Booking, error) {
	ctx := context.Background()

	input := &dynamodb.QueryInput{
		TableName:              aws.String(repo.table),
		IndexName:              aws.String("LSI-4"),
		KeyConditionExpression: aws.String("PK = :pk AND #date BETWEEN :startDate AND :endDate"),
		ExpressionAttributeNames: map[string]string{
			"#date": "Date",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":        &types.AttributeValueMemberS{Value: "BOOKING"},
			":startDate": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", startDate)},
			":endDate":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", endDate)},
		},
	}

	result, err := repo.client.Query(ctx, input)
	if err != nil {
		log.Printf("Failed to get bookings by date range: %v", err)
		return nil, fmt.Errorf("failed to get bookings by date range: %w", err)
	}

	if len(result.Items) == 0 {
		return []domain.Booking{}, nil
	}

	var items []dto.BookingDynamoDBItem
	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		log.Printf("Failed to unmarshal bookings: %v", err)
		return nil, fmt.Errorf("failed to unmarshal bookings: %w", err)
	}

	bookings := make([]domain.Booking, len(items))
	for i, item := range items {
		bookings[i] = domain.Booking{
			ID:        item.ID,
			UserID:    item.UserID,
			RoomID:    item.RoomID,
			StartTime: item.StartTime,
			EndTime:   item.EndTime,
			Purpose:   item.Purpose,
			Status:    item.Status,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}

	return bookings, nil
}

func (repo *BookingRepositoryDynamoDB) GetByRoomIDAndDate(roomID string, date int64) ([]domain.Booking, error) {
	ctx := context.Background()

	startOfDay := (date / 86400) * 86400

	input := &dynamodb.QueryInput{
		TableName:              aws.String(repo.table),
		IndexName:              aws.String("LSI-5"),
		KeyConditionExpression: aws.String("PK = :pk AND RoomID = :roomId"),
		FilterExpression:       aws.String("#date = :date"),
		ExpressionAttributeNames: map[string]string{
			"#date": "Date",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: "BOOKING"},
			":roomId": &types.AttributeValueMemberS{Value: roomID},
			":date":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", startOfDay)},
		},
	}

	result, err := repo.client.Query(ctx, input)
	if err != nil {
		log.Printf("Failed to get bookings by room and date: %v", err)
		return nil, fmt.Errorf("failed to get bookings by room and date: %w", err)
	}

	if len(result.Items) == 0 {
		return []domain.Booking{}, nil
	}

	var items []dto.BookingDynamoDBItem
	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		log.Printf("Failed to unmarshal bookings: %v", err)
		return nil, fmt.Errorf("failed to unmarshal bookings: %w", err)
	}

	bookings := make([]domain.Booking, len(items))
	for i, item := range items {
		bookings[i] = domain.Booking{
			ID:        item.ID,
			UserID:    item.UserID,
			RoomID:    item.RoomID,
			StartTime: item.StartTime,
			EndTime:   item.EndTime,
			Purpose:   item.Purpose,
			Status:    item.Status,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}

	return bookings, nil
}
