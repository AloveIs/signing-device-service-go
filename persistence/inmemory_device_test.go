package persistence

// TODO: add other tests to verify the other interface functions

import (
	"sync"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/common"
)

// Test that TransactionalUpdateDevice does not cause race condtions between goroutines
// TODO: add -race option to testing
func TestDeviceTransactionalUpdate(t *testing.T) {
	N := 10000

	db := NewInMemoryDeviceDb()
	wg := &sync.WaitGroup{}

	startDevice := common.DeviceDTO{
		ID:               "1",
		SignatureCounter: 0,
	}

	if err := db.SaveDevice(startDevice); err != nil {
		t.Errorf("Cannot create device: %v", err)
	}

	for i := 0; i < N; i++ {
		wg.Add(1)

		go func(db DeviceRepository) {
			db.TransactionalUpdateDevice(
				startDevice.ID,
				func(device *common.DeviceDTO) error {
					device.SignatureCounter += 1
					return nil
				})
			wg.Done()
		}(db)

	}
	wg.Wait()

	finalDevice, err := db.GetDeviceByID(startDevice.ID)
	if err != nil {
		t.Errorf("Cannot get device: %v", err)
	}

	if finalDevice.SignatureCounter != uint64(N) {
		t.Errorf("Expected %v, got %v", N, finalDevice.SignatureCounter)
	}
}

// TestDevivceCreateListRetrieve verifies the basic CRUD operations for devices:
// 1. Initially confirms the device list is empty
// 2. Creates 3 devices with IDs "1", "2", "3"
// 3. Verifies that listing devices returns all 3 created devices
func TestDevivceCreateListRetrieve(t *testing.T) {

	db := NewInMemoryDeviceDb()

	devices, err := db.ListDevices()

	if err != nil {
		t.Errorf("Cannot list devices: %v", err)
	}

	if len(devices) != 0 {
		t.Errorf("Expected 0 devices, got %v", len(devices))
	}

	idsToCreate := []string{"1", "2", "3"}
	for _, id := range idsToCreate {
		err := db.SaveDevice(common.DeviceDTO{ID: id})
		if err != nil {
			t.Errorf("Cannot create device: %v", err)
		}
	}

	devices, err = db.ListDevices()

	if err != nil {
		t.Errorf("Cannot list devices: %v", err)
	}

	if len(devices) != len(idsToCreate) {
		t.Errorf("Expected %v devices, got %v", len(idsToCreate), len(devices))
	}

}
