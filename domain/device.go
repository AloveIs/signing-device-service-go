package domain

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/google/uuid"
)

type Device struct {
	ID        string `json:"id"`
	Algorithm string `json:"algorithm"`
	// TODO: change this with the correct type
	signer           crypto.Signer
	Label            string `json:"label"`
	signatureCounter uint64
	LastSignature    string
}

func generateDeviceId() string {
	// TODO: check uniqueness of the ID and panic behaviour of the function
	return uuid.NewString()
}

func NewDevice(label string, algorithm string) Device {
	// TODO: validate label
	signer, err := crypto.NewSigner(algorithm)

	if err != nil {
		// TODO: handle this
		panic(err)
	}
	// TODO: validate Algorithm
	return Device{
		ID:               generateDeviceId(),
		Label:            label,
		Algorithm:        algorithm,
		signer:           signer,
		signatureCounter: 0,
	}
}
