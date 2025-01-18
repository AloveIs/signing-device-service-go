package domain

import (
	"encoding/base64"
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/common"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/google/uuid"
)

type Device struct {
	ID        string `json:"id"`
	Algorithm string `json:"algorithm"`
	// TODO: change this with the correct type
	signer           crypto.MarshallableSigner
	Label            string `json:"label"`
	signatureCounter uint64
	LastSignature    string
}

func (d Device) ToSerializable() common.SerializableDevice {

	return common.SerializableDevice{
		ID:        d.ID,
		Algorithm: d.Algorithm,
		Label:     d.Label,
		PublicKey: d.signer.PublicKey(),
	}
}

func DeviceFromDTO(dto common.DeviceDTO) (Device, error) {
	var d Device
	d.ID = dto.ID
	d.Label = dto.Label
	d.Algorithm = dto.Algorithm
	d.signatureCounter = dto.SignatureCounter
	d.LastSignature = dto.LastSignature
	signer, err := crypto.UnmarshalSigner(d.Algorithm, dto.PrivateKey)
	if err != nil {
		// TODO: handle this
		panic(err)
	}
	d.signer = signer
	return d, nil
}

func (d Device) ToDTO() common.DeviceDTO {
	publicKey, privateKey, err := d.signer.Marshal()

	if err != nil {
		// TODO: handle this
		panic(err)
	}

	return common.DeviceDTO{
		ID:               d.ID,
		Label:            d.Label,
		Algorithm:        d.Algorithm,
		PrivateKey:       privateKey,
		PublicKey:        publicKey,
		SignatureCounter: d.signatureCounter,
		LastSignature:    d.LastSignature,
	}
}

func generateDeviceId() string {
	// TODO: investigate uniqueness of the ID and panic behaviour of the function
	return uuid.NewString()
}

func NewDevice(label string, algorithm string) (Device, error) {
	if len(label) == 0 {
		return Device{}, fmt.Errorf("label cannot be empty")
	}

	signer, err := crypto.NewSigner(algorithm)

	if err != nil {
		return Device{}, fmt.Errorf("invalid algorithm value")
	}
	return Device{
		ID:               generateDeviceId(),
		Label:            label,
		Algorithm:        algorithm,
		signer:           signer,
		signatureCounter: 0,
	}, nil
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
