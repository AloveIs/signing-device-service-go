package domain_test

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/AloveIs/signing-device-service-go/common"
	"github.com/AloveIs/signing-device-service-go/domain"
	"github.com/AloveIs/signing-device-service-go/persistence"
)

// Helper struct to break down and validate signature components
type signatureDestructed struct {
	counter    int    // Sequential counter for the signature
	dataToSign string // Original data that was signed
	sign       string // The actual signature
	prevSign   string // Previous signature in the chain
}

// createTestServiceInstance creates a new in-memory database and device service for testing
func createTestServiceInstance() *domain.DeviceService {
	deviceDb := persistence.NewInMemoryDeviceDb()
	signatureDb := persistence.NewInMemorySignatureDb()
	return domain.NewDeviceService(deviceDb, signatureDb)
}

// TestCRUDOperations verifies basic Create, Read operations for devices
func TestCRUDOperations(t *testing.T) {
	deviceService := createTestServiceInstance()

	// Verify initial state has no devices
	devices, err := deviceService.GetAllDevices()
	if err != nil {
		t.Errorf("Error listing devices: %v", err)
	}
	if len(devices) != 0 {
		t.Errorf("Expected empty device list, found %d devices", len(devices))
	}

	// Test device creation
	createdDevice, err := deviceService.CreateDevice("RSA", nil)
	if err != nil {
		t.Errorf("Error creating device: %v", err)
	}

	// Verify device retrieval
	retrievedDevice, err := deviceService.GetDeviceByID(createdDevice.ID)
	if err != nil {
		t.Errorf("Error retrieving device: %v", err)
	}

	if retrievedDevice != createdDevice {
		t.Errorf("Retrieved device does not match created device")
	}
}

// TestSignature verifies the signing functionality works correctly
// and proper errors are returned for invalid devices
func TestSignature(t *testing.T) {
	deviceService := createTestServiceInstance()

	// Create a test device
	createdDevice, err := deviceService.CreateDevice("RSA", nil)
	if err != nil {
		t.Errorf("Error creating device: %v", err)
	}

	// Test successful signing
	_, err = deviceService.SignMessageWithDevice(createdDevice.ID, []byte("data"))
	if err != nil {
		t.Errorf("Error signing data: %v", err)
	}

	// Test signing with invalid device ID
	impossibleDeviceID := "####"
	_, err = deviceService.SignMessageWithDevice(impossibleDeviceID, []byte("data"))
	if err == nil {
		t.Error("Expected error for invalid device ID, got nil")
	}
	if !errors.Is(err, domain.ErrDeviceNotFound) {
		t.Errorf("Expected ErrDeviceNotFound, got: %v", err)
	}
}

// TestConcurrentUsers verifies that signatures remain consistent
// when multiple users are signing messages simultaneously
func TestConcurrentUsers(t *testing.T) {
	deviceService := createTestServiceInstance()
	N := 10000

	// Create channel to collect signature results
	channel := make(chan common.Signature, N)

	// Create test device
	createdDevice, err := deviceService.CreateDevice("RSA", nil)
	if err != nil {
		t.Errorf("Error creating device: %v", err)
	}

	// Launch N concurrent signing operations
	wg := &sync.WaitGroup{}
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			signature, err := deviceService.SignMessageWithDevice(createdDevice.ID, []byte("data"))
			if err != nil {
				t.Errorf("Error signing data: %v", err)
			}
			channel <- signature
		}()
	}

	wg.Wait()
	close(channel)

	// Collect and validate all signatures
	signature_results := make(map[int]signatureDestructed)
	for sign := range channel {
		splits := strings.Split(sign.SignedData, "_")
		if len(splits) != 3 {
			t.Errorf("Invalid signature format, expected 3 parts, got %d", len(splits))
		}

		counter, err := strconv.ParseInt(splits[0], 10, 32)
		if err != nil {
			t.Errorf("Error parsing signature counter: %v", err)
		}
		signature_results[int(counter)] = signatureDestructed{
			counter:    int(counter),
			dataToSign: splits[1],
			sign:       sign.Signature,
			prevSign:   splits[2],
		}
	}

	// Verify signature chain integrity
	validateSignatureChain(t, signature_results, N)
}

// TestSequentialSign verifies that signatures remain consistent
// when signing messages sequentially
func TestSequentialSign(t *testing.T) {
	deviceService := createTestServiceInstance()
	N := 10000

	createdDevice, err := deviceService.CreateDevice("RSA", nil)
	if err != nil {
		t.Errorf("Error creating device: %v", err)
	}

	// Generate signatures sequentially
	signature_results := make(map[int]signatureDestructed)
	for i := 0; i < N; i++ {
		signature, err := deviceService.SignMessageWithDevice(createdDevice.ID, []byte("data"))
		if err != nil {
			t.Errorf("Error signing data: %v", err)
		}

		splits := strings.Split(signature.SignedData, "_")
		if len(splits) != 3 {
			t.Errorf("Invalid signature format, expected 3 parts, got %d", len(splits))
		}

		counter, err := strconv.ParseInt(splits[0], 10, 32)
		if err != nil {
			t.Errorf("Error parsing signature counter: %v", err)
		}
		signature_results[int(counter)] = signatureDestructed{
			counter:    int(counter),
			dataToSign: splits[1],
			sign:       signature.Signature,
			prevSign:   splits[2],
		}
	}

	// Verify signature chain integrity
	validateSignatureChain(t, signature_results, N)
}

// validateSignatureChain is a helper function to verify the integrity
// of a chain of signatures
func validateSignatureChain(t *testing.T, signature_results map[int]signatureDestructed, N int) {
	for i := 0; i < N; i++ {
		sign, has := signature_results[i]
		if !has {
			t.Errorf("Missing signature for counter %d", i)
		}

		// Skip first signature as it has no previous signature
		if i == 0 {
			continue
		}

		// Verify signature chain links
		prevsig := signature_results[i-1]
		if prevsig.sign != sign.prevSign {
			t.Errorf("Signature chain broken at %d: previous signature %s doesn't match stored previous %s",
				i, prevsig.sign, sign.prevSign)
		}

		if signature_results[i].counter != i {
			t.Errorf("Invalid counter value at position %d: expected %d, got %d",
				i, i, signature_results[i].counter)
		}
	}
}
