package domain

import (
	"errors"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/common"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

type DeviceService struct {
	repo persistence.DeviceRepository
}

func NewDeviceService(repository persistence.DeviceRepository) *DeviceService {
	return &DeviceService{
		repo: repository,
	}
}

func (s *DeviceService) CreateDevice(label string, algorithm string) (common.SerializableDevice, error) {

	device, err := NewDevice(label, algorithm)
	if err != nil {
		return common.SerializableDevice{}, err
	}
	err = s.repo.CreateDevice(device.ToDTO())
	if err != nil {
		return common.SerializableDevice{}, err
	}
	return device.ToSerializable(), nil
}

func (s *DeviceService) GetAllDevices() ([]common.SerializableDevice, error) {
	DTOdevices, err := s.repo.ListDevices()
	if err != nil {
		return nil, err
	}
	result := make([]common.SerializableDevice, 0, len(DTOdevices))
	for _, DTOd := range DTOdevices {
		d, err := DeviceFromDTO(DTOd)
		if err != nil {
			return nil, err
		}
		result = append(result, d.ToSerializable())
	}

	return result, nil
}

func (s *DeviceService) GetDeviceByID(deviceID string) (common.SerializableDevice, error) {

	deviceDTO, err := s.repo.GetDeviceByID(deviceID)
	if err != nil && errors.Is(err, persistence.ErrDeviceNotFound) {
		return common.SerializableDevice{}, ErrDeviceNotFound
	} else if err != nil {
		return common.SerializableDevice{}, err
	}
	device, err := DeviceFromDTO(deviceDTO)
	if err != nil && errors.Is(err, persistence.ErrDeviceNotFound) {
		return common.SerializableDevice{}, ErrDeviceNotFound
	}
	return device.ToSerializable(), nil
}

type SignedMessageDigest struct {
	Signature  string
	SignedData string
}

func (s *DeviceService) SignMessageWithDevice(deviceID string, message []byte) (SignedMessageDigest, error) {
	// TODO: make the signature result capture more elegant, e.g. add a result interface{} as second argument of updateFn
	var signature string
	var signedData string
	err := s.repo.TransactionalUpdateDevice(deviceID, func(deviceDTO *common.DeviceDTO) error {
		var err error
		device, err := DeviceFromDTO(*deviceDTO)
		if err != nil {
			return err
		}
		signature, signedData, err = device.Sign(message)

		*deviceDTO = device.ToDTO()
		if err != nil {
			return err
		}
		return nil
	})
	if errors.Is(err, persistence.ErrDeviceNotFound) {
		return SignedMessageDigest{}, ErrDeviceNotFound
	}

	return SignedMessageDigest{
		Signature:  signature,
		SignedData: signedData,
	}, nil
}
