package common

type SerializableDevice struct {
	ID        string  `json:"id"`
	Algorithm string  `json:"algorithm"`
	Label     *string `json:"label"`
	PublicKey string  `json:"public_key"`
}

type DeviceDTO struct {
	ID               string
	Label            *string
	Algorithm        string
	PrivateKey       []byte
	PublicKey        []byte
	SignatureCounter uint64
	LastSignature    string
}
