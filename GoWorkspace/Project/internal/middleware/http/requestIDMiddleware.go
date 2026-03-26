package middlewareHTTP

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// contextKey — приватный тип для ключей context, чтобы избежать конфликтов
type contextKey string

// requestIDKey — ключ для хранения request ID в context
const requestIDKey contextKey = "request-id"

// RequestIDMiddleware добавляет или генерирует request ID и кладёт его в context и response header
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Получаем X-Request-ID из заголовка или генерируем новый UUID
		requestID := r.Header.Get("x-request-id")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		// Кладём requestID в context и создаём новый запрос с этим context
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		r = r.WithContext(ctx)

		// Отправляем клиенту заголовок ответа с id
		w.Header().Set("X-Request-ID", requestID)

		// Вызываем следующий handler с обновлённым request
		next.ServeHTTP(w, r)
	})
}

// GetRequestID возвращает request ID из context или пустую строку
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}
