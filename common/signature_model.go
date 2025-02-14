package common

// Signature represents a signature created using a device
// It is meant to be serialized toe xternal services
type Signature struct {
	ID         string `json:"id"`
	DeviceID   string `json:"device_id"`
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}

// SignatureDTO for communicating with the persistence layer
type SignatureDTO struct {
	ID         string
	DeviceID   string
	Signature  string
	SignedData string
}

// ToSignature converts a SignatureDTO to a Signature
func (dto *SignatureDTO) ToSignature() Signature {
	return Signature{
		ID:         dto.ID,
		DeviceID:   dto.DeviceID,
		Signature:  dto.Signature,
		SignedData: dto.SignedData,
	}
}
