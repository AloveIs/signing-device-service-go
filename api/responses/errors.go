package responses

import (
	"fmt"
	"net/http"
)

// ErrorResponse is the generic error API response container.
type APIError struct {
	StatusCode int `json:"status_code"`
	Errors     any `json:"errors"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("Api error: %d", e.StatusCode)
}

func InvalidJSON() APIError {
	return APIError{
		StatusCode: http.StatusBadRequest,
		Errors:     "Invalid JSON request data",
	}
}

func InvalidRequestData(errors map[string]string) APIError {
	return APIError{
		StatusCode: http.StatusUnprocessableEntity,
		Errors:     errors,
	}
}

func UrlNotFoundError() error {
	return APIError{
		StatusCode: http.StatusNotFound,
		Errors:     "Url not found",
	}
}

func NewAPIError(code int, errors any) error {
	return &APIError{
		StatusCode: code,
		Errors:     errors,
	}
}
