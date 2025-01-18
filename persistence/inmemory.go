package persistence

import (
	"fmt"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/common"
)

type InMmemoryDb struct {
	// RWMutex to emulate atomicity of the database
	rwmutex sync.RWMutex
	// Storage method is a map deviceID:device
	db map[string]common.DeviceDTO
}

func (imdb *InMmemoryDb) GetDeviceByID(deviceID string) (common.DeviceDTO, error) {

	imdb.rwmutex.RLock()
	defer imdb.rwmutex.RUnlock()

	val, has := imdb.db[deviceID]
	if !has {
		return common.DeviceDTO{}, ErrDeviceNotFound
	}
	return val, nil
}

func (imdb *InMmemoryDb) CreateDevice(device common.DeviceDTO) error {
	imdb.rwmutex.Lock()
	defer imdb.rwmutex.Unlock()
	deviceID := device.ID
	// Check key collision
	_, has := imdb.db[deviceID]
	if has {
		return fmt.Errorf("invaid device, a device with the same ID found.")
	}
	imdb.db[deviceID] = device
	return nil
}

func (imdb *InMmemoryDb) UpdateDevice(key string, val common.DeviceDTO) (common.DeviceDTO, error) {
	imdb.rwmutex.Lock()
	defer imdb.rwmutex.Unlock()

	prev_val, has := imdb.db[key]
	if has {
		return common.DeviceDTO{}, ErrDeviceNotFound
	}
	imdb.db[key] = val

	return prev_val, nil
}

func (imdb *InMmemoryDb) TransactionalUpdateDevice(deviceID string, updateFn func(device *common.DeviceDTO) error) error {
	imdb.rwmutex.Lock()
	defer imdb.rwmutex.Unlock()

	device, has := imdb.db[deviceID]
	if !has {
		return ErrDeviceNotFound
	}
	if err := updateFn(&device); err != nil {
		return err
	}
	// perform the database update
	imdb.db[deviceID] = device

	return nil
}

func (imdb *InMmemoryDb) ListDevices() ([]common.DeviceDTO, error) {
	imdb.rwmutex.RLock()
	defer imdb.rwmutex.RUnlock()

	recordsCount := len(imdb.db)
	result := make([]common.DeviceDTO, 0, recordsCount)

	for _, record := range imdb.db {
		result = append(result, record)
	}
	return result, nil
}

func NewInMmemoryDb() DeviceRepository {
	return &InMmemoryDb{
		db: make(map[string]common.DeviceDTO),
	}
}
