package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/amangirdhar210/meeting-room/internal/domain"
)

type BookingRepositorySQLite struct {
	db *sql.DB
}

func NewBookingRepositorySQLite(db *sql.DB) domain.BookingRepository {
	return &BookingRepositorySQLite{db: db}
}

func (r *BookingRepositorySQLite) checkAvailability(roomID int64, startTime, endTime time.Time) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM bookings 
		WHERE room_id = ?
		AND ((start_time < ? AND end_time > ?) OR (start_time < ? AND end_time > ?))
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int
	err := r.db.QueryRowContext(ctx, query, roomID, endTime, startTime, endTime, startTime).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (r *BookingRepositorySQLite) Create(booking *domain.Booking) error {
	if booking == nil {
		return domain.ErrInvalidInput
	}

	available, err := r.checkAvailability(booking.RoomID, booking.StartTime, booking.EndTime)
	if err != nil {
		return err
	}
	if !available {
		return domain.ErrConflict
	}

	query := `
		INSERT INTO bookings (user_id, room_id, start_time, end_time, purpose, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = r.db.ExecContext(ctx, query,
		booking.UserID,
		booking.RoomID,
		booking.StartTime,
		booking.EndTime,
		booking.Purpose,
		time.Now(),
		time.Now(),
	)
	return err
}

func (r *BookingRepositorySQLite) GetByID(id int64) (*domain.Booking, error) {
	query := `
		SELECT id, user_id, room_id, start_time, end_time, purpose, created_at, updated_at
		FROM bookings WHERE id = ?
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var b domain.Booking
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID, &b.UserID, &b.RoomID, &b.StartTime, &b.EndTime, &b.Purpose, &b.CreatedAt, &b.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BookingRepositorySQLite) GetAll() ([]domain.Booking, error) {
	query := `
		SELECT id, user_id, room_id, start_time, end_time, purpose, created_at, updated_at 
		FROM bookings ORDER BY start_time DESC
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		var b domain.Booking
		err := rows.Scan(&b.ID, &b.UserID, &b.RoomID, &b.StartTime, &b.EndTime, &b.Purpose, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	if len(bookings) == 0 {
		return nil, domain.ErrNotFound
	}
	return bookings, nil
}

func (r *BookingRepositorySQLite) GetByRoomAndTime(roomID int64, start, end time.Time) ([]domain.Booking, error) {
	query := `
		SELECT id, user_id, room_id, start_time, end_time, purpose, created_at, updated_at
		FROM bookings
		WHERE room_id = ? AND (
			(start_time < ? AND end_time > ?) OR
			(start_time >= ? AND start_time < ?)
		)
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, roomID, end, start, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		var b domain.Booking
		err := rows.Scan(&b.ID, &b.UserID, &b.RoomID, &b.StartTime, &b.EndTime, &b.Purpose, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	if len(bookings) == 0 {
		return nil, domain.ErrNotFound
	}
	return bookings, nil
}

func (r *BookingRepositorySQLite) GetByRoomID(roomID int64) ([]domain.Booking, error) {
	query := `
		SELECT id, user_id, room_id, start_time, end_time, purpose, created_at, updated_at
		FROM bookings
		WHERE room_id = ?
		ORDER BY start_time ASC
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		var b domain.Booking
		err := rows.Scan(&b.ID, &b.UserID, &b.RoomID, &b.StartTime, &b.EndTime, &b.Purpose, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	if len(bookings) == 0 {
		return nil, domain.ErrNotFound
	}
	return bookings, nil
}

func (r *BookingRepositorySQLite) Cancel(id int64) error {
	query := `DELETE FROM bookings WHERE id = ?`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
