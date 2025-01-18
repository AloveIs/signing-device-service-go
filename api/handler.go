package api

import "net/http"

// Interface for HTTP handlers (see net/http.Handler) that can also return errors.
type ErrorableHttpHandler interface {
	// Function to handle an http request. See net/http.Handler.ServeHTTP
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

// Interface for HTTP handler that are aware of the URL prefix they
// are served on.
type URLPrefixHandler interface {
	// Set the base prefix of the URLs served to the handler
	SetPathPrefix(path string)
}

type RoutedHttpHandler interface {
	ErrorableHttpHandler
	URLPrefixHandler
}
