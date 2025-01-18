package domain

import (
	"encoding/base64"
	"fmt"

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
	// TODO: investigate uniqueness of the ID and panic behaviour of the function
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

func (d *Device) composeDataToBeSigned(dataToSign []byte) string {
	// encode message to b64
	dataToSignB64 := base64.StdEncoding.EncodeToString(dataToSign)

	if d.signatureCounter == 0 {
		deviceIDB64 := base64.StdEncoding.EncodeToString([]byte(d.ID))
		return fmt.Sprintf("%d_%s_%s", d.signatureCounter, dataToSignB64, deviceIDB64)
	}
	return fmt.Sprintf("%d_%s_%s", d.signatureCounter, dataToSignB64, d.LastSignature)
}

func (d *Device) Sign(dataToSign []byte) (string, string, error) {

	securedData := d.composeDataToBeSigned(dataToSign)

	signature, err := d.signer.Sign([]byte(securedData))

	if err != nil {
		return "", "", err
	}

	d.signatureCounter++
	d.LastSignature = base64.StdEncoding.EncodeToString([]byte(signature))

	return d.LastSignature, securedData, nil
}
