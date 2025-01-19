package domain

import (
	"errors"
	"fmt"

	"github.com/AloveIs/signing-device-service-go/crypto"
)

// ErrInvalidAlgorithm is returned when an unsupported signing algorithm is specified
var ErrInvalidAlgorithm = errors.New(fmt.Sprintf("invalid algorithm: value must be of values: %s, %s", crypto.AlgoECDSA, crypto.AlgoRSA))

// UnmarshalSigner creates a signer from a serialized private key using the specified algorithm
// It supports RSA and ECDSA algorithms and returns an error for unsupported algorithms
func unmarshalSigner(algorithm string, privateKey []byte) (crypto.MarshallableSigner, error) {
	switch algorithm {
	case crypto.AlgoRSA:
		return crypto.UnmarshalRSASigner(privateKey)
	case crypto.AlgoECDSA:
		return crypto.UnmarshalECDSASigner(privateKey)
	default:
		return nil, ErrInvalidAlgorithm
	}
}

// NewSigner creates a new signer instance for the specified algorithm
// Generates a new key pair and returns a marshallable signer interface
func newSigner(algorithm string) (crypto.MarshallableSigner, error) {
	switch algorithm {
	case crypto.AlgoRSA:
		return crypto.NewRSASigner()
	case crypto.AlgoECDSA:
		return crypto.NewECDSASigner()
	default:
		return nil, NewValidationError([]string{ErrInvalidAlgorithm.Error()})
	}
}
