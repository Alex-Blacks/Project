package transport

import (
	"Goworkspace/Project/domain"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func DecodeJSONBody(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func WriteJSON(w http.ResponseWriter, r *http.Request, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("[ERROR]: %s: %v", r.URL.Path, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func WriteError(w http.ResponseWriter, r *http.Request, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(ErrorResponse{Error: msg}); err != nil {
		log.Printf("[ERROR]: %s: %v", r.URL.Path, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

}

func HelperError(w http.ResponseWriter, r *http.Request, err error, id ...int) {
	status, strState := MapDomainErrorToHTTP(err)
	if len(id) > 0 {
		log.Printf("[ERROR]: %s %s id=%d: %v", r.Method, r.URL.Path, id[0], err)
	} else {
		log.Printf("[ERROR]: %s %s: %v", r.Method, r.URL.Path, err)
	}
	WriteError(w, r, status, strState)
}
func MapDomainErrorToHTTP(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrEmptyName),
		errors.Is(err, domain.ErrBadRequest),
		errors.Is(err, domain.ErrInvalidValue):
		return http.StatusBadRequest, "bad request"
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, "not found"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
