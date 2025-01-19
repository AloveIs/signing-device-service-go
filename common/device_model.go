package common

// Data representing a signature device.
// It is meant to be serialized toe xternal services
type Device struct {
	ID        string  `json:"id"`
	Algorithm string  `json:"algorithm"`
	Label     *string `json:"label"`
	// TODO: check if public key needs to be sent to the client for local verification
	PublicKey string `json:"-"`
}

// DeviceDTO for the device for communicating with the persistence layer
type DeviceDTO struct {
	ID               string
	Label            *string
	Algorithm        string
	PrivateKey       []byte
	PublicKey        []byte
	SignatureCounter uint64
	LastSignature    string
}
