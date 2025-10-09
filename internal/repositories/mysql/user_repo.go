package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/amangirdhar210/meeting-room/internal/domain"
)

type UserRepositoryMySQL struct {
	db *sql.DB
}

func NewUserRepositoryMySQL(db *sql.DB) domain.UserRepository {
	return &UserRepositoryMySQL{db: db}
}

func (r *UserRepositoryMySQL) Create(user *domain.User) error {
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

func (r *UserRepositoryMySQL) FindByEmail(email string) (*domain.User, error) {
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

func (r *UserRepositoryMySQL) GetByID(id int64) (*domain.User, error) {
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

func (r *UserRepositoryMySQL) GetAll() ([]domain.User, error) {
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
