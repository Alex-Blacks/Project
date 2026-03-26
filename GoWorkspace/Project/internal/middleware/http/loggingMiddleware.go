package middlewareHTTP

import (
	"Goworkspace/internal/logging"
	"log/slog"
	"net/http"
	"time"
)

// responseWriter оборачивает http.ResponseWriter для сохранения HTTP-статуса
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader сохраняет статус-код и передаёт вызов оригинальному ResponseWriter
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware логирует HTTP-запросы: старт, завершение, статус и длительность.
// Logger кладётся в context для использования в downstream handlers.
func LoggingMiddleware(next http.Handler, baseLogger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Получаем requestID из context; если нет, используем "unknown"
		requestID := GetRequestID(r.Context())
		if requestID == "" {
			requestID = "unknown"
		}

		// Создаём структурированный logger для middleware и downstream
		logger := baseLogger.With(
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
		)

		// Добавляем logger в context для downstream
		ctx := logging.WithLogger(r.Context(), logger)
		r = r.WithContext(ctx)

		// Оборачиваем ResponseWriter для сохранения статуса
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Логируем старт запроса и фиксируем время начала
		start := time.Now()
		logger.Info("HTTP: request started")

		// Вызываем handler, который пишет ответ через wrapped ResponseWriter
		next.ServeHTTP(wrapped, r)

		// Логируем результат запроса, статус и длительность
		duration := time.Since(start)
		if wrapped.statusCode >= 400 {
			logger.Error("HTTP: request failed",
				"status", wrapped.statusCode,
				"duration", duration,
			)
		} else {
			logger.Info("HTTP: request succeeded",
				"status", wrapped.statusCode,
				"duration", duration,
			)
		}
	})
}
