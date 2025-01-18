package domain

import (
	"encoding/base64"
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/common"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/google/uuid"
)

// Structure representing a signature device and its business logic
type signatureDevice struct {
	ID string
	// TODO: change this name
	signer           crypto.MarshallableSigner
	Label            *string
	signatureCounter uint64
	LastSignature    string
}

func (d *signatureDevice) ToSerializable() common.Device {

	return common.Device{
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

func deviceFromDTO(dto common.DeviceDTO) (signatureDevice, error) {
	var d signatureDevice
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

func (d signatureDevice) ToDTO() common.DeviceDTO {
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

func newDevice(algorithm string, label *string) (signatureDevice, error) {
	signer, err := NewSigner(algorithm)

	if err != nil {
		return signatureDevice{}, fmt.Errorf("invalid algorithm value")
	}
	return signatureDevice{
		ID:               generateDeviceId(),
		Label:            copyString(label),
		signer:           signer,
		signatureCounter: 0,
	}, nil
}

func (d *signatureDevice) composeDataToBeSigned(dataToSign []byte) string {
	// encode message to b64
	dataToSignB64 := base64.StdEncoding.EncodeToString(dataToSign)

	if d.signatureCounter == 0 {
		deviceIDB64 := base64.StdEncoding.EncodeToString([]byte(d.ID))
		return fmt.Sprintf("%d_%s_%s", d.signatureCounter, dataToSignB64, deviceIDB64)
	}
	return fmt.Sprintf("%d_%s_%s", d.signatureCounter, dataToSignB64, d.LastSignature)
}

func (d *signatureDevice) Sign(dataToSign []byte) (string, string, error) {

	securedData := d.composeDataToBeSigned(dataToSign)

	signature, err := d.signer.Sign([]byte(securedData))

	if err != nil {
		return "", "", err
	}

	d.signatureCounter++
	d.LastSignature = base64.StdEncoding.EncodeToString([]byte(signature))

	return d.LastSignature, securedData, nil
}
