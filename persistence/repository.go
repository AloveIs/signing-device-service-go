package persistence

import "github.com/fiskaly/coding-challenges/signing-service-challenge/domain"

// DeviceRepository handles CRUD operations for devices
type DeviceRepository interface {
	// CreateDevice adds a new device
	CreateDevice(device domain.Device) error

	// GetDeviceByID fetches a device by ID
	// Returns ErrDeviceNotFound if the device is not found
	GetDeviceByID(id string) (domain.Device, error)

	// TransactionalUpdateDevice modifies a device within a transaction
	TransactionalUpdateDevice(id string, updateFn func(device *domain.Device) error) error

	// ListDevices returns all devices
	ListDevices() ([]domain.Device, error)
}
