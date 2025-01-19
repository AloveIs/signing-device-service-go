package domain

import (
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
)

func TestSignatureDeviceSignatureCounter(t *testing.T) {
	// Create original device
	signDevice, err := newDevice("RSA", nil)
	if err != nil {
		t.Fatalf("Failed to create device: %v", err)
	}

	initialCounter := signDevice.signatureCounter

	if initialCounter != 0 {
		t.Errorf("Expected signature counter to be 0, got %d", signDevice.signatureCounter)
	}
	// number of signs
	N := 10
	for i := 0; i < N; i++ {

		// Sign something to change the counter
		_, _, err = signDevice.Sign([]byte("test message"))
		if err != nil {
			t.Fatalf("Failed to sign message: %v", err)
		}
		expected_counter := (uint64(i) + initialCounter + 1)
		// check the counter
		if signDevice.signatureCounter != expected_counter {
			t.Errorf("Expected signature counter to be %d, got %d", expected_counter, signDevice.signatureCounter)
		}
	}
}

func TestSignatureDeviceAlgorithms(t *testing.T) {
	algorithms := []string{crypto.AlgoRSA, crypto.AlgoECDSA}

	for _, algo := range algorithms {
		t.Run(algo, func(t *testing.T) {
			// Create device with algorithm
			device, err := newDevice(algo, nil)
			if err != nil {
				t.Fatalf("Failed to create device with %s algorithm: %v", algo, err)
			}

			// Verify algorithm was set correctly
			if device.signer.GetAlgorithm() != algo {
				t.Errorf("Expected algorithm to be %s, got %s", algo, device.signer.GetAlgorithm())
			}

			// Test signing
			_, _, err = device.Sign([]byte("test message"))
			if err != nil {
				t.Errorf("Failed to sign with %s algorithm: %v", algo, err)
			}
		})
	}

	// Test invalid algorithm
	_, err := newDevice("INVALID", nil)
	if err == nil {
		t.Error("Expected error when creating device with invalid algorithm")
	}
}
