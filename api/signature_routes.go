package api

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api/responses"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

var signatureIDPattern = regexp.MustCompile("^([^/]+)$")

func (h *SignatureAPIHandler) SetPathPrefix(prefix string) {
	h.Prefix = prefix
}

func (handler *SignatureAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return handler.RouteRequest(w, r)
}

// Route the http request to the correct handler.
func (handler *SignatureAPIHandler) RouteRequest(w http.ResponseWriter, r *http.Request) error {
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
	// GET /{signatureID}
	case r.Method == http.MethodGet && signatureIDPattern.MatchString(relative):
		signatureID := signatureIDPattern.FindStringSubmatch(relative)[1]
		return handler.Retrieve(signatureID, w, r)
	default:
		return responses.UrlNotFoundError()
	}
}

func (handler *SignatureAPIHandler) Retrieve(signatureID string, w http.ResponseWriter, r *http.Request) error {
	signature, err := handler.service.GetSignatureByID(signatureID)
	if err != nil && errors.Is(err, domain.ErrSignatureNotFound) {
		return responses.NewAPIError(http.StatusNotFound, fmt.Sprintf("device %s not found", signatureID))
	} else if err != nil {
		return err
	}
	WriteAPIResponse(w, http.StatusOK, signature)
	return nil
}

func (handler *SignatureAPIHandler) List(w http.ResponseWriter, r *http.Request) error {
	signatures, err := handler.service.ListSignatures()
	if err != nil {
		return err
	}
	WriteAPIResponse(w, http.StatusOK, signatures)
	return nil
}
