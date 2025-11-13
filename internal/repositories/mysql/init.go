package mysql

import (
	"database/sql"
	"log"
)

const schema = `
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  role TEXT DEFAULT 'user',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rooms (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  room_number INTEGER NOT NULL,
  capacity INTEGER NOT NULL,
  floor INTEGER NOT NULL,
  amenities TEXT, -- store as JSON array string
  status TEXT NOT NULL DEFAULT 'Available', -- "Available" or "In Use"
  location TEXT,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS bookings (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  room_id INTEGER NOT NULL,
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
		_, err := db.Exec(
			`INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)`,
			"Admin", "admin@example.com", "admin123", "admin",
		)
		if err != nil {
			log.Printf("Failed to seed admin user: %v", err)
			return err
		}
		log.Println("Seeded admin user with email: admin@example.com, password: admin123")
	}
	return nil
}
