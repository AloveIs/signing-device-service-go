package api

import "net/http"

// Interface for HTTP handlers that can return errors
type ErrorableHttpHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}
