package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/amangirdhar210/meeting-room/internal/domain"
)

type RoomRepositoryMySQL struct {
	db *sql.DB
}

func NewRoomRepositoryMySQL(db *sql.DB) domain.RoomRepository {
	return &RoomRepositoryMySQL{db: db}
}

func (r *RoomRepositoryMySQL) Create(room *domain.Room) error {
	if room == nil {
		return domain.ErrInvalidInput
	}

	now := time.Now()
	if room.CreatedAt.IsZero() {
		room.CreatedAt = now
	}
	room.UpdatedAt = now

	query := `
		INSERT INTO rooms (name, capacity, location, available, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query,
		room.Name,
		room.Capacity,
		room.Location,
		room.IsAvailable,
		room.CreatedAt,
		room.UpdatedAt,
	)
	return err
}

func (r *RoomRepositoryMySQL) GetAll() ([]domain.Room, error) {
	query := `SELECT id, name, capacity, location, available, created_at, updated_at FROM rooms`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []domain.Room
	for rows.Next() {
		var rm domain.Room
		err := rows.Scan(&rm.ID, &rm.Name, &rm.Capacity, &rm.Location, &rm.IsAvailable, &rm.CreatedAt, &rm.UpdatedAt)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, rm)
	}
	if len(rooms) == 0 {
		return nil, domain.ErrNotFound
	}
	return rooms, nil
}

func (r *RoomRepositoryMySQL) GetByID(id int64) (*domain.Room, error) {
	query := `SELECT id, name, capacity, location, available, created_at, updated_at FROM rooms WHERE id = ?`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rm domain.Room
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rm.ID, &rm.Name, &rm.Capacity, &rm.Location, &rm.IsAvailable, &rm.CreatedAt, &rm.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &rm, nil
}

func (r *RoomRepositoryMySQL) UpdateAvailability(id int64, available bool) error {
	query := `UPDATE rooms SET available = ?, updated_at = ? WHERE id = ?`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, available, time.Now(), id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
