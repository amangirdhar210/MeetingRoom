package service

import (
	"strings"
	"time"

	"github.com/amangirdhar210/meeting-room/internal/domain"
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

	room.CreatedAt = time.Now()
	room.UpdatedAt = time.Now()

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

func (s *RoomServiceImpl) GetRoomByID(id int64) (*domain.Room, error) {
	if id <= 0 {
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
func (s *RoomServiceImpl) DeleteRoomByID(id int64) error {
	if id <= 0 {
		return domain.ErrInvalidInput
	}
	return s.repo.DeleteByID(id)
}

func (s *RoomServiceImpl) SearchRooms(minCapacity, maxCapacity int, floor *int, amenities []string, startTime, endTime *time.Time) ([]domain.Room, error) {
	rooms, err := s.repo.SearchWithFilters(minCapacity, maxCapacity, floor, amenities)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (s *RoomServiceImpl) CheckAvailability(roomID int64, startTime, endTime time.Time) (bool, []domain.Booking, error) {
	if roomID <= 0 {
		return false, nil, domain.ErrInvalidInput
	}

	_, err := s.repo.GetByID(roomID)
	if err != nil {
		return false, nil, err
	}

	return true, []domain.Booking{}, nil
}

func (s *RoomServiceImpl) GetAvailableSlots(roomID int64, date time.Time, slotDuration int) ([]domain.TimeSlot, error) {
	if roomID <= 0 || slotDuration <= 0 {
		return nil, domain.ErrInvalidInput
	}

	_, err := s.repo.GetByID(roomID)
	if err != nil {
		return nil, err
	}

	return []domain.TimeSlot{}, nil
}
