package transport

import (
	"Goworkspace/Project/domain"
)

type CreateRequest struct {
	Name string `json:"name"`
}

type ResponseResult struct {
	Item   *domain.Item `json:"item,omitempty"`
	Status string       `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
