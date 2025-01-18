package domain

import (
	"errors"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/common"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

// Serivce exposing all the business logic operations regarding device managment
// and signing.
type DeviceService struct {
	deviceRepo    persistence.DeviceRepository
	signatureRepo persistence.SignatureRepository
}

func NewDeviceService(devices persistence.DeviceRepository, signatures persistence.SignatureRepository) *DeviceService {
	return &DeviceService{
		deviceRepo:    devices,
		signatureRepo: signatures,
	}
}

// CreateDevice creates a new device with the specified signing algorithm and optional label.
// Returns the created device or an error if the creation fails. If the input values are not
// wrong a ValidationError is returned.
func (s *DeviceService) CreateDevice(algorithm string, label *string) (common.Device, error) {

	device, err := newDevice(algorithm, label)
	if err != nil {
		return common.Device{}, err
	}
	err = s.deviceRepo.CreateDevice(device.ToDTO())
	if err != nil {
		return common.Device{}, err
	}
	return device.ToSerializable(), nil
}

func (s *DeviceService) GetAllDevices() ([]common.Device, error) {
	DTOdevices, err := s.deviceRepo.ListDevices()
	if err != nil {
		return nil, err
	}
	result := make([]common.Device, 0, len(DTOdevices))
	for _, DTOd := range DTOdevices {
		d, err := deviceFromDTO(DTOd)
		if err != nil {
			return nil, err
		}
		result = append(result, d.ToSerializable())
	}

	return result, nil
}

// GetDeviceByID retrieves the device with deviceID.
// If the device is not found ErrDeviceNotFound is returned.
func (s *DeviceService) GetDeviceByID(deviceID string) (common.Device, error) {

	deviceDTO, err := s.deviceRepo.GetDeviceByID(deviceID)
	if err != nil && errors.Is(err, persistence.ErrNotFound) {
		return common.Device{}, ErrDeviceNotFound
	} else if err != nil {
		return common.Device{}, err
	}
	device, err := deviceFromDTO(deviceDTO)
	if err != nil && errors.Is(err, persistence.ErrNotFound) {
		return common.Device{}, ErrDeviceNotFound
	}
	return device.ToSerializable(), nil
}

// SignMessageWithDevice signs a message using the device identified by deviceID.
// Returns the signature and signed data, or ErrDeviceNotFound if the device does not exist.
func (s *DeviceService) SignMessageWithDevice(deviceID string, message []byte) (common.Signature, error) {
	// TODO: make the signature result capture more elegant, e.g. add a result interface{} as second argument of updateFn
	var signatureDTO common.SignatureDTO
	err := s.deviceRepo.TransactionalUpdateDevice(deviceID, func(deviceDTO *common.DeviceDTO) error {
		var err error
		device, err := deviceFromDTO(*deviceDTO)
		if err != nil {
			return err
		}
		signature, signedData, err := device.Sign(message)

		if err != nil {
			return err
		}
		// store the signature
		signatureDTO = common.SignatureDTO{
			DeviceID:   device.ID,
			Signature:  signature,
			SignedData: signedData,
		}
		// TODO this can cause deadlock if the interplay between the two inmemory db gets more complicated (there are 2 independent mutexes)
		err = s.signatureRepo.SaveSignature(signatureDTO)

		if err != nil {
			return err
		}
		// update and device and store it
		*deviceDTO = device.ToDTO()
		if err != nil {
			return err
		}
		return nil
	})
	if errors.Is(err, persistence.ErrNotFound) {
		return common.Signature{}, ErrDeviceNotFound
	}

	return signatureDTO.ToSignature(), nil
}
