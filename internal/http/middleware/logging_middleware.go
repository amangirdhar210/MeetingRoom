package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().Local()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s completed in %s", r.Method, r.RequestURI, time.Since(start))
	})
}
