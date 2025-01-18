package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api/responses"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/service"
)

// DeviceIDPattern matches a device identifier without any path segments (deviceID)
var DeviceIDPattern = regexp.MustCompile("^([^/]+)$")

// DeviceSigningPattern matches a device signing endpoint path (deviceID/sign)
var DeviceSigningPattern = regexp.MustCompile("^([^/]+)/sign$")

type DeviceAPIHandler struct {
	service *service.DeviceService
	Prefix  string
}

func NewDeviceAPIHandler(service *service.DeviceService) *DeviceAPIHandler {
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
	// POST /create
	case r.Method == http.MethodPost && relative == "create":
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

func (handler *DeviceAPIHandler) Retrieve(deviceID string, w http.ResponseWriter, r *http.Request) error {
	device, err := handler.service.GetDeviceByID(deviceID)
	if err != nil && errors.Is(err, domain.ErrDeviceNotFound) {
		return responses.NewAPIError(http.StatusNotFound, fmt.Errorf("device %s not found", deviceID))
	} else if err != nil {
		return err
	}
	WriteAPIResponse(w, http.StatusOK, device)
	return nil
}

func (handler *DeviceAPIHandler) List(w http.ResponseWriter, r *http.Request) error {
	devices, err := handler.service.GetAllDevices()
	if err != nil {
		return err
	}
	WriteAPIResponse(w, http.StatusOK, devices)
	return nil
}

type CreateDeviceRequest struct {
	Label     string `json:"label"`
	Algorithm string `json:"algorithm"`
}

func (handler *DeviceAPIHandler) Create(w http.ResponseWriter, r *http.Request) error {
	// validate input data to be CreateDeviceRequest
	var req CreateDeviceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	// TODO: add values validation
	if err != nil {
		return responses.InvalidJSON()
	}
	device, err := handler.service.CreateDevice(req.Label, req.Algorithm)

	if err != nil {
		return err
	}
	WriteAPIResponse(w, http.StatusCreated, device)
	return nil
}

type SignMessageRequest struct {
	Message  *string `json:"message"`
	IsBase64 *bool   `json:"isBase64"`
}

func (v *SignMessageRequest) GetMessageBytes() []byte {
	if v.IsBase64 != nil && *v.IsBase64 {
		data, _ := base64.StdEncoding.DecodeString(*v.Message)
		// TODO: check err here!
		return data
	}
	return []byte(*v.Message)
}

func (v *SignMessageRequest) Validate() map[string]string {
	errors := make(map[string]string)
	if v.Message == nil {
		errors["message"] = "value is required"
	}
	if v.Message == nil {
		errors["isBase64"] = "value is required"
	}
	return errors
}

func (handler *DeviceAPIHandler) Sign(deviceID string, w http.ResponseWriter, r *http.Request) error {

	// validate input data to be SignMessageRequest
	var payload SignMessageRequest

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return responses.InvalidJSON()
	}

	if errs := payload.Validate(); len(errs) > 0 {
		return responses.InvalidRequestData(errs)
	}
	digest, err := handler.service.SignMessageWithDevice(deviceID, payload.GetMessageBytes())

	if errors.Is(err, domain.ErrDeviceNotFound) {
		return responses.NewAPIError(http.StatusNotFound, fmt.Errorf("device %s not found", deviceID))
	} else if err != nil {
		return err
	}

	WriteAPIResponse(w, http.StatusOK, digest)
	return nil
}
