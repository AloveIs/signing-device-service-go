package api

import (
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api/responses"
)

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	listenAddress string
	mux           *http.ServeMux
}

// NewServer is a factory to instantiate a new Server. Pass the addess is
// a string pair "address:port"
func NewServer(listenAddress string) *Server {
	return &Server{
		listenAddress: listenAddress,
		mux:           http.NewServeMux(),
	}
}

// AddHandler adds an http handler to handle all subpaths of a URL prefix.
func (s *Server) WithHandler(pathPrefix string, handler RoutedHttpHandler) *Server {
	handler.SetPathPrefix(pathPrefix)
	s.mux.HandleFunc(pathPrefix, errorResponseWrapper(handler))
	return s
}

// Run starts the HTTP server with the registered handlers. This function
// is blocking.
func (s *Server) Run() error {
	return http.ListenAndServe(s.listenAddress, s.mux)
}

// errorResponseWrapper wraps an ErrorableHttpHandler to handle error responses,
// returning APIError with proper formatting or a generic internal error.
func errorResponseWrapper(handler ErrorableHttpHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler.ServeHTTP(w, r); err != nil {
			handleError(w, err)
		}
	}
}

func handleError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *responses.APIError:
		WriteErrorResponse(w, e.StatusCode, e.Errors)
	default:
		// TODO: add logging here to store internal errors
		WriteInternalError(w)
	}
}
