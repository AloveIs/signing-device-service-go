package crypto

import (
	"testing"
)

// TODO: Add verify for signature verification

// Test creation of RSA and ECDSA signer
func TestSignerCreation(t *testing.T) {
	// Test RSA signer creation
	rsaKeyPair, err := (&RSAGenerator{}).Generate()
	if err != nil {
		t.Errorf("Failed to generate RSA key pair: %v", err)
	}
	rsaSigner := RSASigner{RSAKeyPair: rsaKeyPair}
	if rsaSigner.GetAlgorithm() != AlgoRSA {
		t.Errorf("Expected %s algorithm, got %s", AlgoRSA, rsaSigner.GetAlgorithm())
	}

	// Test ECDSA signer creation
	eccKeyPair, err := (&ECCGenerator{}).Generate()
	if err != nil {
		t.Errorf("Failed to generate RSA key pair: %v", err)
	}
	eccSigner := ECCSigner{ECCKeyPair: eccKeyPair}
	if eccSigner.GetAlgorithm() != AlgoECDSA {
		t.Errorf("Expected %s algorithm, got %s", AlgoECDSA, rsaSigner.GetAlgorithm())
	}
}

// Test Marshalling of a (RSA) Signer
//   - create a signer
//   - marshall and unmarshall the signer
//   - verify the signature is the same
func TestSignerMarshalUnmarshalRSA(t *testing.T) {
	// Test RSA signer
	rsaKeyPair, err := (&RSAGenerator{}).Generate()
	if err != nil {
		t.Errorf("Failed to generate RSA key pair: %v", err)
	}
	rsaSigner := &RSASigner{
		RSAKeyPair: rsaKeyPair,
	}

	// Marshal and unmarshal
	_, rsaPrivKey, err := rsaSigner.Marshal()
	if err != nil {
		t.Fatalf("Failed to marshal RSA keys: %v", err)
	}

	rsaMarshaler := NewRSAMarshaler()
	newRSAKeyPair, err := rsaMarshaler.Unmarshal(rsaPrivKey)
	if err != nil {
		t.Fatalf("Failed to unmarshal RSA keys: %v", err)
	}

	if !rsaKeyPair.Private.Equal(newRSAKeyPair.Private) {
		t.Error("RSA keys do not match after marshal/unmarshal")
	}
}

// Test Marshalling of a (ECDSA) Signer
//   - create a signer
//   - marshal and unmarshal the signer
//   - verify the private key is the same
func TestSignerMarshalUnmarshalECDSA(t *testing.T) {

	// Test ECC signer
	eccKeyPair, err := (&ECCGenerator{}).Generate()
	if err != nil {
		t.Errorf("Failed to generate ECC key pair: %v", err)
	}
	eccSigner := &ECCSigner{
		ECCKeyPair: eccKeyPair,
	}

	// Marshal and unmarshal
	_, eccPrivKey, err := eccSigner.Marshal()
	if err != nil {
		t.Fatalf("Failed to marshal ECC keys: %v", err)
	}

	eccMarshaler := NewECCMarshaler()
	newECCKeyPair, err := eccMarshaler.Unmarshal(eccPrivKey)

	if err != nil {
		t.Fatalf("Failed to unmarshal ECC keys: %v", err)
	}

	if !eccKeyPair.Private.Equal(newECCKeyPair.Private) {
		t.Error("ECC keys do not match after marshal/unmarshal")
	}
}
