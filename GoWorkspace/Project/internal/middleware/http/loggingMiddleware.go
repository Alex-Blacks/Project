package middlewareHTTP

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("[INFO]: HTTP: method=%s, path=%s, duration=%s", r.Method, r.URL.Path, duration)
	})
}
