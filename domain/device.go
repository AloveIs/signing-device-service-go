package domain

import (
	"encoding/base64"
	"fmt"

	"github.com/AloveIs/signing-device-service-go/common"
	"github.com/AloveIs/signing-device-service-go/crypto"
	"github.com/google/uuid"
)

// signatureDevice represents a signature device and its business logic
// it is meant to be used only inside this package
type signatureDevice struct {
	// Unique identifier
	ID string
	// signer is the object that can sign a message
	signer crypto.MarshallableSigner
	// label is an optional alternative name for the device
	Label *string
	// counter of the number of signature performed
	signatureCounter uint64
	// last signature performed
	LastSignature string
}

// Create a new signature device from and algorithm and an optional label
func newDevice(algorithm string, label *string) (signatureDevice, error) {
	signer, err := newSigner(algorithm)

	if err != nil {
		return signatureDevice{}, err
	}
	return signatureDevice{
		ID:               generateDeviceId(),
		Label:            copyString(label),
		signer:           signer,
		signatureCounter: 0,
	}, nil
}

// Sign a message and return its signature
// returns the signature, the data signed and an error
func (d *signatureDevice) sign(dataToSign []byte) (string, string, error) {

	securedData := d.composeDataToBeSigned(dataToSign)

	signature, err := d.signer.Sign([]byte(securedData))

	if err != nil {
		return "", "", err
	}

	d.signatureCounter++
	d.LastSignature = base64.StdEncoding.EncodeToString([]byte(signature))

	return d.LastSignature, securedData, nil
}

// Convert a device into a DTO that can be exposed to outside ervices
func (d *signatureDevice) ToSerializable() common.Device {

	return common.Device{
		ID:        d.ID,
		Algorithm: d.signer.GetAlgorithm(),
		Label:     copyString(d.Label),
		PublicKey: d.signer.PublicKey(),
	}
}

// Unmarshal a device from its DTO representation
func deviceFromDTO(dto common.DeviceDTO) (signatureDevice, error) {
	var d signatureDevice
	d.ID = dto.ID
	d.Label = copyString(dto.Label)
	d.signatureCounter = dto.SignatureCounter
	d.LastSignature = dto.LastSignature
	signer, err := unmarshalSigner(dto.Algorithm, dto.PrivateKey)
	if err != nil {
		return d, fmt.Errorf("Error unmarshalling DTO with algorithm %s (device ID %s): %w", dto.Algorithm, dto.ID, err)
	}
	d.signer = signer
	return d, nil
}

// Marshall a deviec into its DTO
func (d signatureDevice) toDTO() common.DeviceDTO {
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

func (d *signatureDevice) composeDataToBeSigned(dataToSign []byte) string {
	// encode message to b64
	dataToSignB64 := base64.StdEncoding.EncodeToString(dataToSign)

	if d.signatureCounter == 0 {
		deviceIDB64 := base64.StdEncoding.EncodeToString([]byte(d.ID))
		return fmt.Sprintf("%d_%s_%s", d.signatureCounter, dataToSignB64, deviceIDB64)
	}
	return fmt.Sprintf("%d_%s_%s", d.signatureCounter, dataToSignB64, d.LastSignature)
}

func copyString(s *string) *string {
	if s == nil {
		return nil
	}
	c := *s
	return &c
}
