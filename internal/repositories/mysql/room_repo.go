package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/amangirdhar210/meeting-room/internal/domain"
)

type RoomRepositorySQLite struct {
	db *sql.DB
}

func NewRoomRepositorySQLite(db *sql.DB) domain.RoomRepository {
	return &RoomRepositorySQLite{db: db}
}

func (r *RoomRepositorySQLite) scanRoom(rows *sql.Rows) (domain.Room, error) {
	var room domain.Room
	var amenitiesJSON string
	err := rows.Scan(&room.ID, &room.Name, &room.RoomNumber, &room.Capacity, &room.Floor, &amenitiesJSON, &room.Status, &room.Location, &room.Description, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}
	if err := json.Unmarshal([]byte(amenitiesJSON), &room.Amenities); err != nil {
		room.Amenities = []string{}
	}
	return room, nil
}

func (r *RoomRepositorySQLite) scanRooms(rows *sql.Rows) ([]domain.Room, error) {
	var rooms []domain.Room
	for rows.Next() {
		room, err := r.scanRoom(rows)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (r *RoomRepositorySQLite) Create(room *domain.Room) error {
	if room == nil {
		return domain.ErrInvalidInput
	}

	now := time.Now().Unix()
	if room.CreatedAt == 0 {
		room.CreatedAt = now
	}
	room.UpdatedAt = now

	amenitiesJson, err := json.Marshal(room.Amenities)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO rooms (id, name, room_number, capacity, floor, amenities, status, location, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, execErr := r.db.ExecContext(ctx, query,
		room.ID,
		room.Name,
		room.RoomNumber,
		room.Capacity,
		room.Floor,
		string(amenitiesJson),
		room.Status,
		room.Location,
		room.Description,
		room.CreatedAt,
		room.UpdatedAt,
	)
	return execErr
}

func (r *RoomRepositorySQLite) GetAll() ([]domain.Room, error) {
	query := `SELECT id, name, room_number, capacity, floor, amenities, status, location, description, created_at, updated_at FROM rooms`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rooms, err := r.scanRooms(rows)
	if err != nil {
		return nil, err
	}
	if len(rooms) == 0 {
		return nil, domain.ErrNotFound
	}
	return rooms, nil
}

func (r *RoomRepositorySQLite) GetByID(roomID string) (*domain.Room, error) {
	query := `SELECT id, name, room_number, capacity, floor, amenities, status, location, description, created_at, updated_at FROM rooms WHERE id = ?`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room domain.Room
	var amenitiesJSON string
	err := r.db.QueryRowContext(ctx, query, roomID).Scan(
		&room.ID, &room.Name, &room.RoomNumber, &room.Capacity, &room.Floor, &amenitiesJSON, &room.Status, &room.Location, &room.Description, &room.CreatedAt, &room.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(amenitiesJSON), &room.Amenities); err != nil {
		room.Amenities = []string{}
	}
	return &room, nil
}

func (r *RoomRepositorySQLite) UpdateAvailability(roomID string, roomStatus string) error {
	query := `UPDATE rooms SET status = ?, updated_at = ? WHERE id = ?`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, roomStatus, time.Now().Unix(), roomID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *RoomRepositorySQLite) DeleteByID(roomID string) error {
	if roomID == "" {
		return domain.ErrInvalidInput
	}

	query := `DELETE FROM rooms WHERE id = ?`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, roomID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *RoomRepositorySQLite) SearchWithFilters(minCapacity, maxCapacity int, floorNumber *int, requestedAmenities []string) ([]domain.Room, error) {
	query := `SELECT id, name, room_number, capacity, floor, amenities, status, location, description, created_at, updated_at FROM rooms WHERE 1=1`
	queryArgs := []any{}

	if minCapacity > 0 {
		query += ` AND capacity >= ?`
		queryArgs = append(queryArgs, minCapacity)
	}
	if maxCapacity > 0 {
		query += ` AND capacity <= ?`
		queryArgs = append(queryArgs, maxCapacity)
	}
	if floorNumber != nil {
		query += ` AND floor = ?`
		queryArgs = append(queryArgs, *floorNumber)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []domain.Room
	for rows.Next() {
		room, err := r.scanRoom(rows)
		if err != nil {
			return nil, err
		}

		if len(requestedAmenities) > 0 {
			hasAllAmenities := true
			for _, requestedAmenity := range requestedAmenities {
				found := false
				for _, roomAmenity := range room.Amenities {
					if roomAmenity == requestedAmenity {
						found = true
						break
					}
				}
				if !found {
					hasAllAmenities = false
					break
				}
			}
			if !hasAllAmenities {
				continue
			}
		}

		rooms = append(rooms, room)
	}

	if len(rooms) == 0 {
		return []domain.Room{}, nil
	}
	return rooms, nil
}
