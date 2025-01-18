package api

import "net/http"

// Interface for HTTP handlers (see net/http.Handler) that can also return errors.
type ErrorableHttpHandler interface {
	// Function to handle an http request. See net/http.Handler.ServeHTTP
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

type URLPrefixHandler interface {
	// Let the handler be aware of the possible URL prefixes
	SetPathPrefix(path string)
}

type RoutedHttpHandler interface {
	ErrorableHttpHandler
	URLPrefixHandler
}
