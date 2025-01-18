package domain

import (
	"errors"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
)

var ErrInvalidAlgorithm = errors.New("Invalid algorithm: value must be either RSA or ECC")

func UnmarshalSigner(algorithm string, privateKey []byte) (crypto.MarshallableSigner, error) {
	switch algorithm {
	case "RSA":
		g := crypto.NewRSAMarshaler()
		keys, err := g.Unmarshal(privateKey)
		return &crypto.RSASigner{RSAKeyPair: keys}, err
	case "ECC":
		g := crypto.NewECCMarshaler()
		keys, err := g.Unmarshal(privateKey)
		return &crypto.ECCSigner{ECCKeyPair: keys}, err
	default:
		return nil, ErrInvalidAlgorithm
	}
}

func NewSigner(algorithm string) (crypto.MarshallableSigner, error) {
	switch algorithm {
	case "RSA":
		g := &crypto.RSAGenerator{}
		keys, err := g.Generate()
		return &crypto.RSASigner{RSAKeyPair: keys}, err
	case "ECC":
		g := &crypto.ECCGenerator{}
		keys, err := g.Generate()
		return &crypto.ECCSigner{ECCKeyPair: keys}, err
	default:
		return nil, ErrInvalidAlgorithm
	}
}
