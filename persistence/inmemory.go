package persistence

import (
	"fmt"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

type InMmemoryDb struct {
	// RWMutex to emulate atomicity of the database
	rwmutex sync.RWMutex
	// Storage method is a map deviceID:device
	db map[string]domain.Device
}

func (imdb *InMmemoryDb) GetDeviceByID(deviceID string) (domain.Device, error) {

	imdb.rwmutex.RLock()
	defer imdb.rwmutex.RUnlock()

	val, has := imdb.db[deviceID]
	if !has {
		return domain.Device{}, ErrDeviceNotFound
	}
	return val, nil
}

func (imdb *InMmemoryDb) CreateDevice(device domain.Device) error {
	imdb.rwmutex.Lock()
	defer imdb.rwmutex.Unlock()
	deviceID := device.ID
	// TODO: should I check for key collisions?
	_, has := imdb.db[deviceID]
	if has {
		return fmt.Errorf("invaid device, a device with the same ID found.")
	}
	imdb.db[deviceID] = device
	return nil
}

func (imdb *InMmemoryDb) UpdateDevice(key string, val domain.Device) (domain.Device, error) {
	imdb.rwmutex.Lock()
	defer imdb.rwmutex.Unlock()

	// TODO: check if previous must be returned
	prev_val, has := imdb.db[key]
	if has {
		return domain.Device{}, ErrDeviceNotFound
	}
	imdb.db[key] = val

	return prev_val, nil
}

func (imdb *InMmemoryDb) TransactionalUpdateDevice(deviceID string, updateFn func(device *domain.Device) error) error {
	imdb.rwmutex.Lock()
	defer imdb.rwmutex.Unlock()

	// TODO: check if previous must be returned
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

func (imdb *InMmemoryDb) ListDevices() ([]domain.Device, error) {
	imdb.rwmutex.RLock()
	defer imdb.rwmutex.RUnlock()

	recordsCount := len(imdb.db)
	result := make([]domain.Device, 0, recordsCount)

	for _, record := range imdb.db {
		result = append(result, record)
	}
	return result, nil
}

func NewInMmemoryDb() DeviceRepository {
	return &InMmemoryDb{
		db: make(map[string]domain.Device),
	}
}
