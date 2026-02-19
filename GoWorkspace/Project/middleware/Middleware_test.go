package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRecovery_middleware(t *testing.T) {
	handlerPanic := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { panic("Boom") })

	recoveryHandler := RecoveryMiddleware(handlerPanic)

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()

	recoveryHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("Expected code 500, got: %d", rec.Code)
	}

	if rec.Body.String() == "" {
		t.Fatalf("Expected body, got: %s", rec.Body.String())
	}

}

func TestTimeout_middleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { <-r.Context().Done() })

	timeoutHandler := TimeoutMiddleware(1 * time.Second)(handler)

	req := httptest.NewRequest(http.MethodGet, "/slow", nil)
	rec := httptest.NewRecorder()

	timeoutHandler.ServeHTTP(rec, req)

	http.TimeoutHandler(timeoutHandler, 1*time.Second, "timeout")
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected code 200, got: %d", rec.Code)
	}

}
