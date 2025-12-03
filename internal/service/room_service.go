package service

import (
	"strings"
	"time"

	"github.com/amangirdhar210/meeting-room/internal/domain"
	"github.com/google/uuid"
)

type RoomServiceImpl struct {
	repo domain.RoomRepository
}

func NewRoomService(repo domain.RoomRepository) domain.RoomService {
	return &RoomServiceImpl{repo: repo}
}

func (s *RoomServiceImpl) AddRoom(room *domain.Room) error {
	if room == nil {
		return domain.ErrInvalidInput
	}

	room.Name = strings.TrimSpace(room.Name)
	room.Location = strings.TrimSpace(room.Location)

	if room.Name == "" || room.Capacity <= 0 || room.Location == "" || room.RoomNumber <= 0 || room.Floor < 0 {
		return domain.ErrInvalidInput
	}
	if room.Status == "" {
		room.Status = "Available"
	}
	if room.Amenities == nil {
		room.Amenities = []string{}
	}

	room.ID = uuid.New().String()
	room.CreatedAt = time.Now().Unix()
	room.UpdatedAt = time.Now().Unix()

	return s.repo.Create(room)
}

func (s *RoomServiceImpl) GetAllRooms() ([]domain.Room, error) {
	rooms, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	if len(rooms) == 0 {
		return nil, domain.ErrNotFound
	}
	return rooms, nil
}

func (s *RoomServiceImpl) GetRoomByID(id string) (*domain.Room, error) {
	if id == "" {
		return nil, domain.ErrInvalidInput
	}
	room, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, domain.ErrNotFound
	}
	return room, nil
}
func (s *RoomServiceImpl) DeleteRoomByID(id string) error {
	if id == "" {
		return domain.ErrInvalidInput
	}
	return s.repo.DeleteByID(id)
}

func (s *RoomServiceImpl) SearchRooms(minCapacity, maxCapacity int, floor *int, amenities []string, startTime, endTime *int64) ([]domain.Room, error) {
	rooms, err := s.repo.SearchWithFilters(minCapacity, maxCapacity, floor, amenities)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (s *RoomServiceImpl) CheckAvailability(roomID string, startTime, endTime int64) (bool, []domain.Booking, error) {
	if roomID == "" {
		return false, nil, domain.ErrInvalidInput
	}

	_, err := s.repo.GetByID(roomID)
	if err != nil {
		return false, nil, err
	}

	return true, []domain.Booking{}, nil
}

func (s *RoomServiceImpl) GetAvailableSlots(roomID string, date int64, slotDuration int) ([]domain.TimeSlot, error) {
	if roomID == "" || slotDuration <= 0 {
		return nil, domain.ErrInvalidInput
	}

	_, err := s.repo.GetByID(roomID)
	if err != nil {
		return nil, err
	}

	return []domain.TimeSlot{}, nil
}
