package repository

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const schema = `
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  role TEXT DEFAULT 'user',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rooms (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  room_number INTEGER NOT NULL,
  capacity INTEGER NOT NULL,
  floor INTEGER NOT NULL,
  amenities TEXT,
  status TEXT NOT NULL DEFAULT 'Available',
  location TEXT,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS bookings (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  room_id TEXT NOT NULL,
  start_time DATETIME NOT NULL,
  end_time DATETIME NOT NULL,
  purpose TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (room_id) REFERENCES rooms(id)
);
`

func InitSQLite(db *sql.DB) error {
	_, err := db.Exec(schema)
	if err != nil {
		return err
	}

	var count int
	row := db.QueryRow(`SELECT COUNT(*) FROM users WHERE email = ?`, "admin@example.com")
	if err := row.Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash admin password: %v", err)
			return err
		}

		adminID := uuid.New().String()
		_, err = db.Exec(
			`INSERT INTO users (id, name, email, password, role) VALUES (?, ?, ?, ?, ?)`,
			adminID, "Admin", "admin@example.com", string(hashedPassword), "admin",
		)
		if err != nil {
			log.Printf("Failed to seed admin user: %v", err)
			return err
		}
		log.Println("Seeded admin user with email: admin@example.com, password: admin123")
	}
	return nil
}
