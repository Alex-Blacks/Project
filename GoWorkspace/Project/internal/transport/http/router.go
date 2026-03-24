package transport

import (
	middleware "Goworkspace/internal/middleware/http"
	"Goworkspace/internal/service"
	"time"

	"github.com/go-chi/chi/v5"
)

func NewRouter(service *service.Service) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RecoveryMiddleware)
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.TimeoutMiddleware(60 * time.Second))

	r.Post("/item", PostHandler(service))
	r.Get("/item/{id}", GetHandler(service))
	r.Delete("/item/{id}", DeleteHandler(service))

	return r
}
