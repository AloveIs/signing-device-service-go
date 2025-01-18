package persistence

// TODO: add other tests to verify the other interface functions

import (
	"sync"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/common"
)

// Test that TransactionalUpdateDevice does not cause race condtions between goroutines
// TODO: add -race option to testing
func TestTransactionalUpdate(t *testing.T) {
	N := 10000

	db := NewInMmemoryDb()
	wg := &sync.WaitGroup{}

	startDevice := common.DeviceDTO{
		ID:               "1",
		SignatureCounter: 0,
	}

	if err := db.CreateDevice(startDevice); err != nil {
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
