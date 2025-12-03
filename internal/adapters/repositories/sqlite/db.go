package repository

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type DBConfig struct {
	Path string
}

func NewSQLiteConnection(cfg DBConfig) (*sql.DB, error) {
	dir := filepath.Dir(cfg.Path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create db directory: %w", err)
		}
	}
	db, err := sql.Open("sqlite3", cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite connection: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping SQLite: %w", err)
	}
	return db, nil
}

func CloseDB(db *sql.DB) error {
	if db == nil {
		return nil
	}
	return db.Close()
}
