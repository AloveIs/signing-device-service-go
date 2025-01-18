package service

import (
	"errors"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
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

func (s *DeviceService) CreateDevice(label string, algorithm string) (domain.Device, error) {

	device, err := domain.NewDevice(label, algorithm)
	if err != nil {
		return device, err
	}
	err = s.repo.CreateDevice(device)
	if err != nil {
		return domain.Device{}, err
	}
	return device, nil
}

func (s *DeviceService) GetAllDevices() ([]domain.Device, error) {
	return s.repo.ListDevices()
}

func (s *DeviceService) GetDeviceByID(deviceID string) (domain.Device, error) {

	device, err := s.repo.GetDeviceByID(deviceID)
	if errors.Is(err, persistence.ErrDeviceNotFound) {
		return device, domain.ErrDeviceNotFound
	}
	return device, nil
}

type SignedMessageDigest struct {
	Signature  string
	SignedData string
}

func (s *DeviceService) SignMessageWithDevice(deviceID string, message []byte) (SignedMessageDigest, error) {
	// TODO: make the result capture more elegant, e.g. add a result interface{} as second argument of updateFn
	var signature string
	var signedData string
	err := s.repo.TransactionalUpdateDevice(deviceID, func(device *domain.Device) error {
		var err error
		signature, signedData, err = device.Sign(message)

		if err != nil {
			return err
		}
		return nil
	})
	if errors.Is(err, persistence.ErrDeviceNotFound) {
		return SignedMessageDigest{}, domain.ErrDeviceNotFound
	}

	return SignedMessageDigest{
		Signature:  signature,
		SignedData: signedData,
	}, nil
}
