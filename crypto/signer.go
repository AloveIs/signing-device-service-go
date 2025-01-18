package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
	PublicKey() string
	//PrivateKey() string
}

type MarshallableSigner interface {
	Signer
	Marshal() ([]byte, []byte, error)
	//Unmarshal(data []byte) (Signer, error)
}

// computeHash calculates the SHA-256 hash of the input data
func computeHash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func (keys *ECCKeyPair) Sign(dataToBeSigned []byte) ([]byte, error) {
	hash := computeHash(dataToBeSigned)

	signature, err := ecdsa.SignASN1(rand.Reader, keys.Private, hash[:])
	if err != nil {
		return []byte{}, fmt.Errorf("failed to sign data with ECC: %v", err)
	}
	return signature, nil
}

func (keys *RSAKeyPair) Sign(dataToBeSigned []byte) ([]byte, error) {

	hash := computeHash(dataToBeSigned)
	signature, err := rsa.SignPKCS1v15(rand.Reader, keys.Private, crypto.SHA256, hash[:])
	if err != nil {
		return []byte{}, fmt.Errorf("failed to sign data with RSA: %v", err)
	}

	return signature, nil
}

func UnmarshalSigner(algorithm string, privateKey []byte) (MarshallableSigner, error) {
	switch algorithm {
	case "RSA":
		g := NewRSAMarshaler()
		keys, err := g.Unmarshal(privateKey)
		return &RSASigner{RSAKeyPair: keys}, err
	case "ECC":
		g := NewECCMarshaler()
		keys, err := g.Unmarshal(privateKey)
		return &ECCSigner{ECCKeyPair: keys}, err
	default:
		return nil, ErrUnsupportedAlgorithm
	}
}

func NewSigner(algorithm string) (MarshallableSigner, error) {
	switch algorithm {
	case "RSA":
		g := &RSAGenerator{}
		keys, err := g.Generate()
		return &RSASigner{RSAKeyPair: keys}, err
	case "ECC":
		g := &ECCGenerator{}
		keys, err := g.Generate()
		return &ECCSigner{ECCKeyPair: keys}, err
	default:
		return nil, ErrUnsupportedAlgorithm
	}
}
