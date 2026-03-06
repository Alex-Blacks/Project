package transport

import (
	"Goworkspace/Project/domain"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func PostHandler(src *domain.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req CreateRequest

		if err := DecodeJSONBody(r, &req); err != nil {
			HelperError(w, r, err)
			return
		}

		item, err := src.Create(r.Context(), req.Name)
		if err != nil {
			HelperError(w, r, err)
			return
		}

		res := &ResponseResult{Item: &item, Status: "Create OK"}
		WriteJSON(w, r, http.StatusCreated, res)

		log.Printf("[INFO]: %s %s: successful: id=%d", r.Method, r.URL.Path, item.ID)
	})
}

func GetHandler(src *domain.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		strID := chi.URLParam(r, "id")
		reqID, err := strconv.Atoi(strID)
		if err != nil || reqID < 1 {
			HelperError(w, r, domain.ErrInvalidValue)
			return
		}

		item, err := src.Get(r.Context(), reqID)
		if err != nil {
			HelperError(w, r, err, reqID)
			return
		}

		res := ResponseResult{Item: &item, Status: "Get OK"}
		WriteJSON(w, r, http.StatusOK, res)

		log.Printf("[INFO]: %s %s: successful: id=%d", r.Method, r.URL.Path, reqID)
	})
}

func DeleteHandler(src *domain.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		strID := chi.URLParam(r, "id")
		reqID, err := strconv.Atoi(strID)
		if err != nil || reqID < 1 {
			HelperError(w, r, domain.ErrInvalidValue)
			return
		}

		if err := src.Delete(r.Context(), reqID); err != nil {
			HelperError(w, r, err, reqID)
			return
		}

		res := ResponseResult{Status: "Delete OK"}
		WriteJSON(w, r, http.StatusOK, res)

		log.Printf("[INFO]: %s %s: successful: id=%d", r.Method, r.URL.Path, reqID)
	})
}
