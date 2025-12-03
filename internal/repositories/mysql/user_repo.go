package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/amangirdhar210/meeting-room/internal/domain"
)

type UserRepositorySQLite struct {
	db *sql.DB
}

func NewUserRepositorySQLite(db *sql.DB) domain.UserRepository {
	return &UserRepositorySQLite{db: db}
}

func (r *UserRepositorySQLite) Create(user *domain.User) error {
	if user == nil {
		return domain.ErrInvalidInput
	}

	query := `
		INSERT INTO users (id, name, email, password, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Name, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *UserRepositorySQLite) FindByEmail(userEmail string) (*domain.User, error) {
	query := `SELECT id, name, email, password, role, created_at, updated_at FROM users WHERE email = ? LIMIT 1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, userEmail).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositorySQLite) GetByID(userID string) (*domain.User, error) {
	query := `SELECT id, name, email, password, role, created_at, updated_at FROM users WHERE id = ?`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositorySQLite) GetAll() ([]domain.User, error) {
	query := `SELECT id, name, email, role, created_at, updated_at FROM users`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if len(users) == 0 {
		return nil, domain.ErrNotFound
	}
	return users, nil
}

func (r *UserRepositorySQLite) DeleteByID(userID string) error {
	query := `DELETE FROM users WHERE id = ?`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, userID)
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
