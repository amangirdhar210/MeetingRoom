package dynamodb

import (
	"github.com/amangirdhar210/meeting-room/internal/core/domain"
	"github.com/amangirdhar210/meeting-room/internal/core/ports"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
	return nil
}
func (repo *BookingRepositoryDynamoDB) GetByID(id string) (*domain.Booking, error) {
	return nil, nil
}
func (repo *BookingRepositoryDynamoDB) GetAll() ([]domain.Booking, error) {
	return nil, nil
}
func (repo *BookingRepositoryDynamoDB) GetByRoomAndTime(roomID string, start, end int64) ([]domain.Booking, error) {
	return nil, nil
}
func (repo *BookingRepositoryDynamoDB) GetByRoomID(roomID string) ([]domain.Booking, error) {
	return nil, nil
}
func (repo *BookingRepositoryDynamoDB) GetByUserID(userID string) ([]domain.Booking, error) {
	return nil, nil
}
func (repo *BookingRepositoryDynamoDB) Cancel(id string) error {
	return nil
}
func (repo *BookingRepositoryDynamoDB) GetByDateRange(startDate, endDate int64) ([]domain.Booking, error) {
	return nil, nil
}
func (repo *BookingRepositoryDynamoDB) GetByRoomIDAndDate(roomID string, date int64) ([]domain.Booking, error) {
	return nil, nil
}

