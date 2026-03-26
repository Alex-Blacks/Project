package middlewareHTTP

import (
	"Goworkspace/internal/logging"
	"encoding/json"
	"net/http"
)

// RecoveryMiddleware защищает сервер от паники в обработчиках.
// Если происходит panic, логирует её и возвращает клиенту HTTP 500 с JSON-ответом.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// defer-функция ловит панику после выполнения handler
		defer func() {
			// Получаем logger из context
			logger := logging.LoggerFromContext(r.Context())

			// Проверяем, была ли паника
			if rec := recover(); rec != nil {
				// Логируем panic с подробностями
				logger.Error("HTTP:[RECOVERY]: panic",
					"recovered", rec,
				)

				// Формируем ответ клиенту с кодом 500
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				// Пытаемся вернуть JSON с сообщением об ошибке
				if err := json.NewEncoder(w).Encode(map[string]string{"Error": "internal server error"}); err != nil {
					// Если не удалось сформировать ответ, логируем вторичную ошибку
					logger.Error("failed to encode recovery response",
						"error", err,
					)
				}
			}
		}()

		// Вызываем следующий handler
		next.ServeHTTP(w, r)
	})
}
