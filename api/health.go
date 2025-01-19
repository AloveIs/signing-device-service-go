package api

import (
	"net/http"

	"github.com/AloveIs/signing-device-service-go/api/responses"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
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

// Ignore the prefix, needed to implement the RoutedHttpHandler interface
// TODO: remove this and improve the builder pattern on Server to accept different interfaces
func (h *HealthHandler) SetPathPrefix(path string) {}
