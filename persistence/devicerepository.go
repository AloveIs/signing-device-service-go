package persistence

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/common"
)

// DeviceRepository handles CRUD operations for devices
type DeviceRepository interface {
	// CreateDevice adds a new device to the repository
	SaveDevice(device common.DeviceDTO) error

	// GetDeviceByID fetches a device by ID
	// Returns ErrDeviceNotFound if the device is not found
	GetDeviceByID(id string) (common.DeviceDTO, error)

	// TransactionalUpdateDevice modifies a device within a SQL-like transaction with
	// the provided updateFn function.
	TransactionalUpdateDevice(id string, updateFn func(device *common.DeviceDTO) error) error

	// ListDevices returns all devices
	ListDevices() ([]common.DeviceDTO, error)
}
