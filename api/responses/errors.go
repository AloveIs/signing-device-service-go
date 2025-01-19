package responses

import (
	"fmt"
	"net/http"
)

// ErrorResponse is the generic error API response container.
// Handlers can return this as an error to be processed by the middleware
// to return a JSON response to the client.
type APIError struct {
	StatusCode int `json:"status_code"`
	Errors     any `json:"errors"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("Api error: %d", e.StatusCode)
}

func NewAPIError(code int, errors any) error {
	return &APIError{
		StatusCode: code,
		Errors:     errors,
	}
}

// Create a APIError representing an (syntactic) invalid JSON
func InvalidJSON() *APIError {
	return &APIError{
		StatusCode: http.StatusBadRequest,
		Errors:     "Invalid JSON request data",
	}
}

func InvalidRequestData(errors []string) *APIError {
	return &APIError{
		StatusCode: http.StatusUnprocessableEntity,
		Errors:     errors,
	}
}

func UrlNotFoundError() error {
	return &APIError{
		StatusCode: http.StatusNotFound,
		Errors:     "Url not found",
	}
}
