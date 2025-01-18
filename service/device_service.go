package service

import (
	"errors"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

type DeviceBusinessLogicService struct {
	repo persistence.DeviceRepository
}

func NewDeviceBusinessLogicService(repository persistence.DeviceRepository) *DeviceBusinessLogicService {
	return &DeviceBusinessLogicService{
		repo: repository,
	}
}

func (s *DeviceBusinessLogicService) CreateDevice(label string, algorithm string) (domain.Device, error) {

	device := domain.NewDevice(label, algorithm)
	// TODO: add error for wrong algorithm/valdiation
	err := s.repo.CreateDevice(device)

	if err != nil {
		return domain.Device{}, err
	}
	return device, nil
}

func (s *DeviceBusinessLogicService) GetAllDevices() ([]domain.Device, error) {
	return s.repo.ListDevices()
}

func (s *DeviceBusinessLogicService) GetDeviceByID(deviceID string) (domain.Device, error) {

	return s.repo.GetDeviceByID(deviceID)
}

type SignedMessageDigest struct {
	Signature  string
	SignedData string
}

func NewDeviceService(repository persistence.DeviceRepository) *DeviceBusinessLogicService {
	return &DeviceBusinessLogicService{
		repo: repository,
	}
}

func (s *DeviceBusinessLogicService) SignMessageWithDevice(deviceID string, message []byte) (SignedMessageDigest, error) {
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
	if errors.Is(persistence.ErrDeviceNotFound, err) {
		return SignedMessageDigest{}, domain.ErrDeviceNotFound
	}

	return SignedMessageDigest{
		Signature:  signature,
		SignedData: signedData,
	}, nil
}
