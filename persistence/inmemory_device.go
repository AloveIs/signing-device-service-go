package persistence

import (
	"sync"

	"github.com/AloveIs/signing-device-service-go/common"
)

// InMemoryDeviceDb implements an in-memory database for storing device records
// using a map with read-write mutex for concurrent access control
type InMemoryDeviceDb struct {
	// RWMutex to emulate atomicity of the database
	rwmutex sync.RWMutex
	// Storage method is a map deviceID:device
	db map[string]common.DeviceDTO
}

func (imdb *InMemoryDeviceDb) GetDeviceByID(deviceID string) (common.DeviceDTO, error) {

	imdb.rwmutex.RLock()
	defer imdb.rwmutex.RUnlock()

	val, has := imdb.db[deviceID]
	if !has {
		return common.DeviceDTO{}, ErrNotFound
	}
	return val, nil
}

func (imdb *InMemoryDeviceDb) SaveDevice(device common.DeviceDTO) error {
	imdb.rwmutex.Lock()
	defer imdb.rwmutex.Unlock()
	deviceID := device.ID
	// Check key collision
	_, has := imdb.db[deviceID]
	if has {
		return ErrIdKeyCollision
	}
	imdb.db[deviceID] = device
	return nil
}

func (imdb *InMemoryDeviceDb) UpdateDevice(key string, val common.DeviceDTO) (common.DeviceDTO, error) {
	imdb.rwmutex.Lock()
	defer imdb.rwmutex.Unlock()

	prev_val, has := imdb.db[key]
	if has {
		return common.DeviceDTO{}, ErrNotFound
	}
	imdb.db[key] = val

	return prev_val, nil
}

func (imdb *InMemoryDeviceDb) TransactionalUpdateDevice(deviceID string, updateFn func(device *common.DeviceDTO) error) error {
	imdb.rwmutex.Lock()
	defer imdb.rwmutex.Unlock()

	device, has := imdb.db[deviceID]
	if !has {
		return ErrNotFound
	}
	if err := updateFn(&device); err != nil {
		return err
	}
	// perform the database update
	imdb.db[deviceID] = device

	return nil
}

func (imdb *InMemoryDeviceDb) ListDevices() ([]common.DeviceDTO, error) {
	imdb.rwmutex.RLock()
	defer imdb.rwmutex.RUnlock()

	recordsCount := len(imdb.db)
	result := make([]common.DeviceDTO, 0, recordsCount)

	for _, record := range imdb.db {
		result = append(result, record)
	}
	return result, nil
}

func NewInMemoryDeviceDb() DeviceRepository {
	return &InMemoryDeviceDb{
		db: make(map[string]common.DeviceDTO),
	}
}
