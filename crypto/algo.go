package crypto

import "crypto/sha256"

const (
	AlgoRSA   = "RSA"
	AlgoECDSA = "ECDSA"
)

// computeHash calculates the SHA-256 hash of the input data
func computeHash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}
