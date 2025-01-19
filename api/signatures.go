package api

import (
	"github.com/AloveIs/signing-device-service-go/domain"
)

// SignatureAPIHandler routes and exposes http requests
// to the device service.
type SignatureAPIHandler struct {
	service *domain.SignatureService
	Prefix  string
}

// Create a new SignatureAPIHandler wrapping the provided service
func NewSignatureAPIHandler(service *domain.SignatureService) *SignatureAPIHandler {
	return &SignatureAPIHandler{
		service: service,
		Prefix:  "",
	}
}
