package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/amangirdhar210/meeting-room/internal/app"
	"github.com/amangirdhar210/meeting-room/internal/repositories/mysql"
)

func main() {
	jwtSecret := "supersecretkey"

	cfg := mysql.DBConfig{
		User:     "root",
		Password: "password",
		Host:     "127.0.0.1",
		Port:     3306,
		Name:     "meeting_room_db",
	}

	db, err := mysql.NewMySQLConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	router := app.SetupRouter(db, jwtSecret)

	addr := ":8080"
	fmt.Printf("Server running on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
