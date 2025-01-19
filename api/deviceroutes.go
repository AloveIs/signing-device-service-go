package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api/responses"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

func (handler *DeviceAPIHandler) Retrieve(deviceID string, w http.ResponseWriter, r *http.Request) error {
	device, err := handler.service.GetDeviceByID(deviceID)
	if err != nil && errors.Is(err, domain.ErrDeviceNotFound) {
		return responses.NewAPIError(http.StatusNotFound, fmt.Sprintf("device %s not found", deviceID))
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
	Label     *string `json:"label"`
	Algorithm string  `json:"algorithm"`
}

func (v *CreateDeviceRequest) Validate() []string {
	errors := make([]string, 0)
	if len(v.Algorithm) == 0 {
		errors = append(errors, "algorithm: value is required")
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

func (handler *DeviceAPIHandler) Create(w http.ResponseWriter, r *http.Request) error {
	// validate input data to be CreateDeviceRequest
	var req CreateDeviceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	// TODO: add values validation
	if err != nil {
		return responses.InvalidJSON()
	}

	errs := req.Validate()

	if len(errs) > 0 {
		return responses.InvalidRequestData(errs)
	}

	device, err := handler.service.CreateDevice(req.Algorithm, req.Label)

	var validationErr *domain.ValidationError = &domain.ValidationError{}

	if errors.As(err, validationErr) {
		return responses.NewAPIError(http.StatusBadRequest, validationErr.Errors)
	} else if err != nil {
		return err
	}
	WriteAPIResponse(w, http.StatusCreated, device)
	return nil
}

type SignMessageRequest struct {
	Message  *string `json:"message"`
	IsBase64 *bool   `json:"isBase64"`
}

func (v *SignMessageRequest) getMessageBytes() []byte {
	// TODO: add mechanism to check input has been validated
	if v.IsBase64 != nil && *v.IsBase64 {
		data, _ := base64.StdEncoding.DecodeString(*v.Message)
		return data
	}
	return []byte(*v.Message)
}

func (v *SignMessageRequest) Validate() ([]byte, []string) {
	errors := make([]string, 0)
	if v.Message == nil {
		errors = append(errors, "message: value is required")
	}
	if v.IsBase64 == nil {
		errors = append(errors, "isBase64: value is required")
	}
	if v.IsBase64 != nil && v.Message != nil && *v.IsBase64 {
		_, err := base64.StdEncoding.DecodeString(*v.Message)
		if err != nil {
			errors = append(errors, "message cannot be decoded into base64")
		}
	}
	if len(errors) > 0 {
		return nil, errors
	}
	return v.getMessageBytes(), nil
}

func (handler *DeviceAPIHandler) Sign(deviceID string, w http.ResponseWriter, r *http.Request) error {

	// validate input data to be SignMessageRequest
	var payload SignMessageRequest

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return responses.InvalidJSON()
	}
	messageBytes, errs := payload.Validate()
	if len(errs) != 0 {
		return responses.InvalidRequestData(errs)
	}
	digest, err := handler.service.SignMessageWithDevice(deviceID, messageBytes)

	if errors.Is(err, domain.ErrDeviceNotFound) {
		return responses.NewAPIError(http.StatusNotFound, fmt.Sprintf("device %s not found", deviceID))
	} else if err != nil {
		return err
	}

	WriteAPIResponse(w, http.StatusCreated, digest)
	return nil
}
