package api

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/AloveIs/signing-device-service-go/api/responses"
	"github.com/AloveIs/signing-device-service-go/domain"
)

// DeviceAPIHandler routes and exposes http requests to the device service.
type DeviceAPIHandler struct {
	service *domain.DeviceService
	Prefix  string
}

// Create a new DeviceAPIHandler using the provided service
func NewDeviceAPIHandler(service *domain.DeviceService) *DeviceAPIHandler {
	return &DeviceAPIHandler{
		service: service,
		Prefix:  "",
	}
}

// RouteRequest routes an http request to its handler.
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
	case r.Method == http.MethodGet && deviceIDPattern.MatchString(relative):
		deviceID := deviceIDPattern.FindStringSubmatch(relative)[1]
		return handler.Retrieve(deviceID, w, r)
	// POST /{deviceID}/sign
	case r.Method == http.MethodPost && deviceSigningPattern.MatchString(relative):
		deviceID := deviceSigningPattern.FindStringSubmatch(relative)[1]
		return handler.Sign(deviceID, w, r)
	default:
		return responses.UrlNotFoundError()
	}
}

// Matches a device identifier without any path segments (deviceID)
var deviceIDPattern = regexp.MustCompile("^([^/]+)$")

// Matches a device signing endpoint path (deviceID/sign)
var deviceSigningPattern = regexp.MustCompile("^([^/]+)/sign$")

func (h *DeviceAPIHandler) SetPathPrefix(prefix string) {
	h.Prefix = prefix
}

func (handler *DeviceAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return handler.RouteRequest(w, r)
}
