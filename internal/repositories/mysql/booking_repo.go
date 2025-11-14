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

func (r *BookingRepositorySQLite) scanBooking(rows *sql.Rows) (domain.Booking, error) {
	var booking domain.Booking
	err := rows.Scan(&booking.ID, &booking.UserID, &booking.RoomID, &booking.StartTime, &booking.EndTime, &booking.Purpose, &booking.CreatedAt, &booking.UpdatedAt)
	return booking, err
}

func (r *BookingRepositorySQLite) scanBookings(rows *sql.Rows) ([]domain.Booking, error) {
	var bookings []domain.Booking
	for rows.Next() {
		booking, err := r.scanBooking(rows)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}
	return bookings, nil
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

	var conflictCount int
	err := r.db.QueryRowContext(ctx, query, roomID, endTime, startTime, endTime, startTime).Scan(&conflictCount)
	if err != nil {
		return false, err
	}
	return conflictCount == 0, nil
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

func (r *BookingRepositorySQLite) GetByID(bookingID int64) (*domain.Booking, error) {
	query := `
		SELECT id, user_id, room_id, start_time, end_time, purpose, created_at, updated_at
		FROM bookings WHERE id = ?
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var booking domain.Booking
	err := r.db.QueryRowContext(ctx, query, bookingID).Scan(
		&booking.ID, &booking.UserID, &booking.RoomID, &booking.StartTime, &booking.EndTime, &booking.Purpose, &booking.CreatedAt, &booking.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &booking, nil
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

	bookings, err := r.scanBookings(rows)
	if err != nil {
		return nil, err
	}
	if len(bookings) == 0 {
		return nil, domain.ErrNotFound
	}
	return bookings, nil
}

func (r *BookingRepositorySQLite) GetByRoomAndTime(roomID int64, startTime, endTime time.Time) ([]domain.Booking, error) {
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

	rows, err := r.db.QueryContext(ctx, query, roomID, endTime, startTime, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings, err := r.scanBookings(rows)
	if err != nil {
		return nil, err
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

	bookings, err := r.scanBookings(rows)
	if err != nil {
		return nil, err
	}
	if len(bookings) == 0 {
		return nil, domain.ErrNotFound
	}
	return bookings, nil
}

func (r *BookingRepositorySQLite) Cancel(bookingID int64) error {
	query := `DELETE FROM bookings WHERE id = ?`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, bookingID)
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

func (r *BookingRepositorySQLite) GetByDateRange(startDate, endDate time.Time) ([]domain.Booking, error) {
	query := `
		SELECT id, user_id, room_id, start_time, end_time, purpose, created_at, updated_at
		FROM bookings
		WHERE start_time >= ? AND end_time <= ?
		ORDER BY start_time ASC
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanBookings(rows)
}

func (r *BookingRepositorySQLite) GetByRoomIDAndDate(roomID int64, targetDate time.Time) ([]domain.Booking, error) {
	startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := `
		SELECT id, user_id, room_id, start_time, end_time, purpose, created_at, updated_at
		FROM bookings
		WHERE room_id = ? AND start_time >= ? AND start_time < ?
		ORDER BY start_time ASC
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, roomID, startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanBookings(rows)
}
