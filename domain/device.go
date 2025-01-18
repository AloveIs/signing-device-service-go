package domain

import (
	"encoding/base64"
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/common"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/google/uuid"
)

type Device struct {
	ID string
	// TODO: remove this and encapsulate it into the signer
	//Algorithm string `json:"algorithm"`
	// TODO: change this with the correct type
	signer           crypto.MarshallableSigner
	Label            *string
	signatureCounter uint64
	LastSignature    string
}

func (d Device) ToSerializable() common.SerializableDevice {

	return common.SerializableDevice{
		ID:        d.ID,
		Algorithm: d.signer.GetAlgorithm(),
		Label:     copyString(d.Label),
		PublicKey: d.signer.PublicKey(),
	}
}

func copyString(s *string) *string {
	if s == nil {
		return nil
	}
	c := *s
	return &c
}

func DeviceFromDTO(dto common.DeviceDTO) (Device, error) {
	var d Device
	d.ID = dto.ID
	d.Label = copyString(dto.Label)
	d.signatureCounter = dto.SignatureCounter
	d.LastSignature = dto.LastSignature
	signer, err := UnmarshalSigner(dto.Algorithm, dto.PrivateKey)
	if err != nil {
		return d, fmt.Errorf("Error unmarshalling DTO with algorithm %s (device ID %s): %w", dto.Algorithm, dto.ID, err)
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
		Algorithm:        d.signer.GetAlgorithm(),
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

func NewDevice(algorithm string, label *string) (Device, error) {
	signer, err := NewSigner(algorithm)

	if err != nil {
		return Device{}, fmt.Errorf("invalid algorithm value")
	}
	return Device{
		ID:               generateDeviceId(),
		Label:            copyString(label),
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
