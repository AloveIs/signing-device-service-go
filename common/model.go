package common

// Data representing a Device that can be serialized to
// external services
type Device struct {
	ID        string  `json:"id"`
	Algorithm string  `json:"algorithm"`
	Label     *string `json:"label"`
	PublicKey string  `json:"public_key"`
}

// DTO for the device for communicating with the persistence layer
type DeviceDTO struct {
	ID               string
	Label            *string
	Algorithm        string
	PrivateKey       []byte
	PublicKey        []byte
	SignatureCounter uint64
	LastSignature    string
}
