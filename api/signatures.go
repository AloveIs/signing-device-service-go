package api

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

// DeviceAPIHandler routes and exposes http requests
// to the device service.
type SignatureAPIHandler struct {
	service *domain.SignatureService
	Prefix  string
}

// Create a new SignatureAPIHandler using the service
func NewSignatureAPIHandler(service *domain.SignatureService) *SignatureAPIHandler {
	return &SignatureAPIHandler{
		service: service,
		Prefix:  "",
	}
}
