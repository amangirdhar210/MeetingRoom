package dynamodb

import (
	"github.com/amangirdhar210/meeting-room/internal/core/domain"
	"github.com/amangirdhar210/meeting-room/internal/core/ports"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
	return nil
}

func (repo *RoomRepositoryDynamoDB) GetAll() ([]domain.Room, error) {
	return nil, nil
}

func (repo *RoomRepositoryDynamoDB) GetByID(id string) (*domain.Room, error) {
	return nil, nil
}

func (repo *RoomRepositoryDynamoDB) UpdateAvailability(id string, status string) error {
	return nil
}

func (repo *RoomRepositoryDynamoDB) DeleteByID(id string) error {
	return nil
}

func (repo *RoomRepositoryDynamoDB) SearchWithFilters(minCapacity, maxCapacity int, floor *int, amenities []string) ([]domain.Room, error) {
	return nil, nil
}
