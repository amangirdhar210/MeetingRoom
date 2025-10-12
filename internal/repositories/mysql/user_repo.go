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
		INSERT INTO users (name, email, password, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query,
		user.Name, user.Email, user.Password, user.Role, time.Now(), time.Now(),
	)
	return err
}

func (r *UserRepositorySQLite) FindByEmail(email string) (*domain.User, error) {
	query := `SELECT id, name, email, password, role, created_at, updated_at FROM users WHERE email = ? LIMIT 1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepositorySQLite) GetByID(id int64) (*domain.User, error) {
	query := `SELECT id, name, email, password, role, created_at, updated_at FROM users WHERE id = ?`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u domain.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
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
		var u domain.User
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if len(users) == 0 {
		return nil, domain.ErrNotFound
	}
	return users, nil
}

func (r *UserRepositorySQLite) DeleteByID(id int64) error {
	query := `DELETE FROM users WHERE id = ?`
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
