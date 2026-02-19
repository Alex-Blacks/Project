package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("[RECOVERY]: %s %s: panic recovered: %v", r.Method, r.URL.Path, rec)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				if err := json.NewEncoder(w).Encode(map[string]string{"Error": "internal server error"}); err != nil {
					log.Printf("[ERROR]: %s %s: %v", r.Method, r.URL.Path, err)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}
