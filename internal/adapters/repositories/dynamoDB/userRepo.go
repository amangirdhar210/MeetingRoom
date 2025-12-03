package dynamodb

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/amangirdhar210/meeting-room/internal/core/domain"
	"github.com/amangirdhar210/meeting-room/internal/core/ports"
)

type UserRepositoryDynamoDB struct {
	client *dynamodb.Client
	table  string
}

func toDomainUser(item map[string]types.AttributeValue) (domain.User, error) {
	user := domain.User{}

	if id, ok := item["ID"].(*types.AttributeValueMemberS); ok {
		user.ID = id.Value
	}
	if name, ok := item["Name"].(*types.AttributeValueMemberS); ok {
		user.Name = name.Value
	}
	if email, ok := item["Email"].(*types.AttributeValueMemberS); ok {
		user.Email = email.Value
	}
	if password, ok := item["Password"].(*types.AttributeValueMemberS); ok {
		user.Password = password.Value
	}
	if role, ok := item["Role"].(*types.AttributeValueMemberS); ok {
		user.Role = role.Value
	}
	if createdAt, ok := item["CreatedAt"].(*types.AttributeValueMemberN); ok {
		if timestamp, err := strconv.ParseInt(createdAt.Value, 10, 64); err == nil {
			user.CreatedAt = timestamp
		}
	}
	if updatedAt, ok := item["UpdatedAt"].(*types.AttributeValueMemberN); ok {
		if timestamp, err := strconv.ParseInt(updatedAt.Value, 10, 64); err == nil {
			user.UpdatedAt = timestamp
		}
	}

	return user, nil
}

func NewUserRepositoryDynamoDB(client *dynamodb.Client, tableName string) ports.UserRepository {
	return &UserRepositoryDynamoDB{
		client: client,
		table:  tableName,
	}
}

func (repo *UserRepositoryDynamoDB) findUserIdByEmail(email string) (string, error) {
	if email == "" {
		return "", domain.ErrInvalidInput
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	input := &dynamodb.QueryInput{
		TableName:              aws.String(repo.table),
		KeyConditionExpression: aws.String("PK = :pk AND SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "USER"},
			":sk": &types.AttributeValueMemberS{Value: email},
		},
	}

	result, err := repo.client.Query(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to query user: %w", err)
	}

	if len(result.Items) == 0 {
		return "", domain.ErrNotFound
	}

	resultItem := result.Items[0]
	if id, ok := resultItem["ID"].(*types.AttributeValueMemberS); ok {
		return id.Value, nil
	}

	return "", domain.ErrNotFound
}

func (repo *UserRepositoryDynamoDB) FindByEmail(email string) (*domain.User, error) {
	userId, err := repo.findUserIdByEmail(email)
	if err != nil {
		return nil, err
	}

	user, err := repo.GetByID(userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *UserRepositoryDynamoDB) GetByID(userID string) (*domain.User, error) {
	if userID == "" {
		return nil, domain.ErrInvalidInput
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	input := &dynamodb.QueryInput{
		TableName:              aws.String(repo.table),
		KeyConditionExpression: aws.String("PK = :pk AND SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "USER"},
			":sk": &types.AttributeValueMemberS{Value: "USER#" + userID},
		},
	}

	result, err := repo.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query user by ID: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, domain.ErrNotFound
	}

	user, err := toDomainUser(result.Items[0])
	if err != nil {
		return nil, fmt.Errorf("failed to convert to domain user: %w", err)
	}

	return &user, nil
}

func (repo *UserRepositoryDynamoDB) Create(user *domain.User) error {
	if user == nil {
		return domain.ErrInvalidInput
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	emailLookupItem := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: "USER"},
		"SK": &types.AttributeValueMemberS{Value: user.Email},
		"ID": &types.AttributeValueMemberS{Value: user.ID},
	}

	userDataItem := map[string]types.AttributeValue{
		"PK":        &types.AttributeValueMemberS{Value: "USER"},
		"SK":        &types.AttributeValueMemberS{Value: "USER#" + user.ID},
		"ID":        &types.AttributeValueMemberS{Value: user.ID},
		"Name":      &types.AttributeValueMemberS{Value: user.Name},
		"Email":     &types.AttributeValueMemberS{Value: user.Email},
		"Password":  &types.AttributeValueMemberS{Value: user.Password},
		"Role":      &types.AttributeValueMemberS{Value: user.Role},
		"CreatedAt": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", user.CreatedAt)},
		"UpdatedAt": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", user.UpdatedAt)},
	}

	transactItems := []types.TransactWriteItem{
		{
			Put: &types.Put{
				TableName: aws.String(repo.table),
				Item:      emailLookupItem,
			},
		},
		{
			Put: &types.Put{
				TableName: aws.String(repo.table),
				Item:      userDataItem,
			},
		},
	}

	_, err := repo.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (repo *UserRepositoryDynamoDB) GetAll() ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	input := &dynamodb.QueryInput{
		TableName:              aws.String(repo.table),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "USER"},
			":sk": &types.AttributeValueMemberS{Value: "USER#"},
		},
	}

	result, err := repo.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query all users: %w", err)
	}

	if len(result.Items) == 0 {
		return []domain.User{}, nil
	}

	users := make([]domain.User, 0, len(result.Items))
	for _, item := range result.Items {
		user, err := toDomainUser(item)
		if err != nil {
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

func (repo *UserRepositoryDynamoDB) DeleteByID(userID string) error {
	if userID == "" {
		return domain.ErrInvalidInput
	}

	user, err := repo.GetByID(userID)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	transactItems := []types.TransactWriteItem{
		{
			Delete: &types.Delete{
				TableName: aws.String(repo.table),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: "USER"},
					"SK": &types.AttributeValueMemberS{Value: user.Email},
				},
			},
		},
		{
			Delete: &types.Delete{
				TableName: aws.String(repo.table),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: "USER"},
					"SK": &types.AttributeValueMemberS{Value: "USER#" + userID},
				},
			},
		},
	}

	_, err = repo.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
