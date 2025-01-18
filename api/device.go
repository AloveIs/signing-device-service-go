package api

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api/responses"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

// Matches a device identifier without any path segments (deviceID)
var DeviceIDPattern = regexp.MustCompile("^([^/]+)$")

// Matches a device signing endpoint path (deviceID/sign)
var DeviceSigningPattern = regexp.MustCompile("^([^/]+)/sign$")

// DeviceAPIHandler routes and exposes http requests
// to the device service.
type DeviceAPIHandler struct {
	service *domain.DeviceService
	Prefix  string
}

// Create a new DeviceAPIHandler using the service
func NewDeviceAPIHandler(service *domain.DeviceService) *DeviceAPIHandler {
	return &DeviceAPIHandler{
		service: service,
		Prefix:  "",
	}
}

func (h *DeviceAPIHandler) SetPathPrefix(prefix string) {
	h.Prefix = prefix
}

func (handler *DeviceAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return handler.RouteRequest(w, r)
}

// Route the http request to the correct handler.
func (handler *DeviceAPIHandler) RouteRequest(w http.ResponseWriter, r *http.Request) error {
	fullpath := r.URL.Path
	relative, found := strings.CutPrefix(fullpath, handler.Prefix)
	if !found {
		// TODO: this is an internal error also (mismatched prefix)
		return responses.UrlNotFoundError()
	}

	switch {
	// GET /
	case r.Method == http.MethodGet && relative == "":
		return handler.List(w, r)
	// POST /
	case r.Method == http.MethodPost && relative == "":
		return handler.Create(w, r)
	// GET /{deviceID}
	case r.Method == http.MethodGet && DeviceIDPattern.MatchString(relative):
		deviceID := DeviceIDPattern.FindStringSubmatch(relative)[1]
		return handler.Retrieve(deviceID, w, r)
	// POST /{deviceID}/sign
	case r.Method == http.MethodPost && DeviceSigningPattern.MatchString(relative):
		deviceID := DeviceSigningPattern.FindStringSubmatch(relative)[1]
		return handler.Sign(deviceID, w, r)
	default:
		return responses.UrlNotFoundError()
	}
}
