package middlewareHTTP

import (
	"context"
	"net/http"
	"time"
)

// TimeoutMiddleware создаёт middleware, который ограничивает время обработки запроса.
// Если handler не завершится за указанный timeout, context автоматически отменяется.
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Создаём context с таймаутом
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel() // гарантируем освобождение ресурсов

			// Передаём новый context в handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
