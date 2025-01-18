package domain

import (
	"errors"
	"fmt"
)

// TODO: make this its own struct to be consistent
var ErrDeviceNotFound = errors.New("Device not found")

type ValidationError struct {
	Errors []string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("Validation errors: %v", e.Errors)
}

func NewValidationError(errors []string) *ValidationError {
	return &ValidationError{Errors: errors}
}
