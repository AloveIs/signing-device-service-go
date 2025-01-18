package api

import (
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api/responses"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

type HealthHandler struct {
}

// Health evaluates the health of the service and writes a standardized response.
func (h *HealthHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) error {
	if request.Method != http.MethodGet {
		return responses.NewAPIError(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
	health := HealthResponse{
		Status:  "pass",
		Version: "v0",
	}
	WriteAPIResponse(response, http.StatusOK, health)
	return nil
}
