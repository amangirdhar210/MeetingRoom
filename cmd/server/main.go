package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/amangirdhar210/meeting-room/internal/app"
	"github.com/amangirdhar210/meeting-room/internal/http/middleware"
	"github.com/amangirdhar210/meeting-room/internal/repositories/mysql"
)

func main() {
	jwtSecret := "supersecretkey"

	cfg := mysql.DBConfig{
		Path: "./meeting_room_db.sqlite",
	}

	db, err := mysql.NewSQLiteConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to SQLite: %v", err)
	}
	defer db.Close()
	if err := mysql.InitSQLite(db); err != nil {
		log.Fatalf("Failed to initialize SQLite schema: %v", err)
	}

	router := app.SetupRouter(db, jwtSecret)
	router = middleware.CORSMiddleware(router)

	addr := ":8080"
	fmt.Printf("Server running on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
