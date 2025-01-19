// Package domain implements the core business logic for the signing service
package domain

import (
	"errors"
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
)

// ErrInvalidAlgorithm is returned when an unsupported signing algorithm is specified
var ErrInvalidAlgorithm = errors.New(fmt.Sprintf("invalid algorithm: value must be of values: %s, %s", crypto.AlgoECDSA, crypto.AlgoRSA))

// UnmarshalSigner creates a signer from a serialized private key using the specified algorithm
// It supports RSA and ECDSA algorithms and returns an error for unsupported algorithms
func UnmarshalSigner(algorithm string, privateKey []byte) (crypto.MarshallableSigner, error) {
	switch algorithm {
	case crypto.AlgoRSA:
		g := crypto.NewRSAMarshaler()
		keys, err := g.Unmarshal(privateKey)
		return &crypto.RSASigner{RSAKeyPair: keys}, err
	case crypto.AlgoECDSA:
		g := crypto.NewECCMarshaler()
		keys, err := g.Unmarshal(privateKey)
		return &crypto.ECCSigner{ECCKeyPair: keys}, err
	default:
		return nil, ErrInvalidAlgorithm
	}
}

// NewSigner creates a new signer instance for the specified algorithm
// Generates a new key pair and returns a marshallable signer interface
func NewSigner(algorithm string) (crypto.MarshallableSigner, error) {
	switch algorithm {
	case crypto.AlgoRSA:
		g := &crypto.RSAGenerator{}
		keys, err := g.Generate()
		return &crypto.RSASigner{RSAKeyPair: keys}, err
	case crypto.AlgoECDSA:
		g := &crypto.ECCGenerator{}
		keys, err := g.Generate()
		return &crypto.ECCSigner{ECCKeyPair: keys}, err
	default:
		return nil, ErrInvalidAlgorithm
	}
}
