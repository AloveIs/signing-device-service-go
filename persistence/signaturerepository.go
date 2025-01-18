package persistence

import "github.com/fiskaly/coding-challenges/signing-service-challenge/common"

type SignatureRepository interface {
	// SaveSignature stores a signature in the repository
	SaveSignature(signature common.SignatureDTO) error

	// GetSignaturesByDeviceID retrieves all signatures for a given device ID
	GetSignaturesByDeviceID(deviceID string) ([]common.SignatureDTO, error)
	// GetSignatureByID retrieves a signature by its ID
	GetSignatureByID(signatureID string) (common.SignatureDTO, error)
	// ListSignatures returns all signatures
	ListSignatures() ([]common.SignatureDTO, error)
}
