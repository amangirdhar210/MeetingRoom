package app

import (
	"log"
	"net/http"
)

func StartServer(addr string, handler http.Handler) {
	log.Printf("Server running on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
