// package mysql

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"time"

// 	_ "github.com/go-sql-driver/mysql"
// )

// type DBConfig struct {
// 	User     string
// 	Password string
// 	Host     string
// 	Port     int
// 	Name     string
// }

// func NewMySQLConnection(cfg DBConfig) (*sql.DB, error) {
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
// 		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
// 	)

// 	db, err := sql.Open("mysql", dsn)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open MySQL connection: %w", err)
// 	}

// 	db.SetMaxOpenConns(25)
// 	db.SetMaxIdleConns(25)
// 	db.SetConnMaxLifetime(5 * time.Minute)

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	if err := db.PingContext(ctx); err != nil {
// 		return nil, fmt.Errorf("failed to ping MySQL: %w", err)
// 	}

// 	return db, nil
// }

// func CloseDB(db *sql.DB) error {
// 	if db == nil {
// 		return nil
// 	}
// 	return db.Close()
// }

package mysql

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
