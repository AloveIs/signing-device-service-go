package api

import (
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api/responses"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	listenAddress string
	repo          persistence.DeviceRepository
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string, deviceRepository persistence.DeviceRepository) *Server {
	return &Server{
		listenAddress: listenAddress,
		repo:          deviceRepository,
	}
}

// errorResponseWrapper takes an ErrorableHttpHandler and returns a standard http.HandlerFunc.
// It wraps the handler to properly format and write any errors that occur during request processing.
// If the handler returns an APIError, it will be formatted accordingly, otherwise a generic internal error is returned.
func errorResponseWrapper(handler ErrorableHttpHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler.ServeHTTP(w, r)
		if err == nil {
			return
		}
		if apiErr, ok := err.(responses.APIError); ok {
			// Write the api error formatted
			WriteErrorResponse(w, apiErr.StatusCode, apiErr.Errors)
		} else {
			// Write a geneal error so not to leak errors information
			WriteInternalError(w)
		}
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	mux := http.NewServeMux()
	healthHandler := &HealthHandler{}
	deviceAPIHandler := NewDeviceAPIHandler("/api/v0/devices/", s.repo)
	// TODO: make it more DRY so not to repeat the path prefix
	mux.HandleFunc("/api/v0/health", errorResponseWrapper(healthHandler))
	mux.HandleFunc(deviceAPIHandler.Prefix, errorResponseWrapper(deviceAPIHandler))

	return http.ListenAndServe(s.listenAddress, mux)
}
