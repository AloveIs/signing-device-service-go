package persistence

import (
	"testing"

	"github.com/AloveIs/signing-device-service-go/common"
)

// TestSignatureCreateListRetrieve tests the core functionality of the signature database:
// - Creating signatures for different devices
// - Listing all signatures
// - Retrieving signatures by device ID
// The test verifies that signatures are correctly stored and can be retrieved both
// globally and filtered by device.
func TestSignatureCreateListRetrieve(t *testing.T) {
	// Initialize a new in-memory signature database
	db := NewInMemorySignatureDb()

	// Initial check - database should be empty
	signatures, err := db.ListSignatures()

	if err != nil {
		t.Errorf("Cannot list signatures: %v", err)
	}
	if len(signatures) != 0 {
		t.Errorf("Expected 0 signatures, got %v", len(signatures))
	}

	// Create test data - three signatures for device A
	idsToCreateDeviceA := []string{"1", "2", "3"}
	for _, id := range idsToCreateDeviceA {
		err := db.SaveSignature(common.SignatureDTO{ID: id, DeviceID: "A"})
		if err != nil {
			t.Errorf("Cannot create device: %v", err)
		}
	}

	// Create test data - two signatures for device B
	idsToCreateDeviceB := []string{"4", "5"}
	for _, id := range idsToCreateDeviceB {
		err := db.SaveSignature(common.SignatureDTO{ID: id, DeviceID: "B"})
		if err != nil {
			t.Errorf("Cannot create device: %v", err)
		}
	}

	// Check total number of signatures across all devices
	expectedTotalSignatures := (len(idsToCreateDeviceA) + len(idsToCreateDeviceB))

	signatures, err = db.ListSignatures()

	if err != nil {
		t.Errorf("Cannot list signatures: %v", err)
	}

	if len(signatures) != expectedTotalSignatures {
		t.Errorf("Expected %v signatures, got %v", expectedTotalSignatures, len(signatures))
	}

	// Verify signatures for device A
	signaturesA, err := db.GetSignaturesByDeviceID("A")
	if err != nil {
		t.Errorf("Cannot list signatures: %v", err)
	}

	if len(signaturesA) != len(idsToCreateDeviceA) {
		t.Errorf("Expected %v signatures, got %v", len(idsToCreateDeviceA), len(signaturesA))
	}

	// Verify signatures for device B
	signaturesB, err := db.GetSignaturesByDeviceID("B")

	if err != nil {
		t.Errorf("Cannot list signatures: %v", err)
	}

	if len(signaturesB) != len(idsToCreateDeviceB) {
		t.Errorf("Expected %v signatures, got %v", len(idsToCreateDeviceB), len(signaturesB))
	}
}
