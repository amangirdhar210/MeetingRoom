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

	if room.Name == "" || room.Capacity <= 0 || room.Location == "" {
		return domain.ErrInvalidInput
	}

	room.CreatedAt = time.Now()
	room.UpdatedAt = time.Now()
	room.IsAvailable = true

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
