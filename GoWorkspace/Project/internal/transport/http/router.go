package transport

import (
	"Goworkspace/internal/logging"
	middleware "Goworkspace/internal/middleware/http"
	"Goworkspace/internal/service"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// NewRouter создаёт chi роутер и подключает все HTTP middleware и маршруты
func NewRouter(svc *service.Service) *chi.Mux {
	r := chi.NewRouter()

	// --- Middleware ---
	// Сначала RequestIDMiddleware добавляет request ID в context и заголовки
	r.Use(middleware.RequestIDMiddleware)

	// LoggingMiddleware логирует start, success/error, duration
	r.Use(func(next http.Handler) http.Handler {
		// обёртка, чтобы передать baseLogger
		return middleware.LoggingMiddleware(next, logging.NewLogger())
	})

	//  Recovery, чтобы поймать паники из всех последующих middleware и handler'ов
	r.Use(middleware.RecoveryMiddleware)

	// TimeoutMiddleware ограничивает время обработки запроса
	r.Use(middleware.TimeoutMiddleware(60 * time.Second))

	// --- Routes ---
	r.Post("/item", PostHandler(svc))
	r.Get("/item/{id}", GetHandler(svc))
	r.Delete("/item/{id}", DeleteHandler(svc))

	return r
}
