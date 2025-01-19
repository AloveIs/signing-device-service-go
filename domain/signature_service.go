package domain

import (
	"errors"

	"github.com/AloveIs/signing-device-service-go/common"
	"github.com/AloveIs/signing-device-service-go/persistence"
)

type SignatureService struct {
	repo persistence.SignatureRepository
}

func NewSignatureService(repository persistence.SignatureRepository) *SignatureService {
	return &SignatureService{
		repo: repository,
	}
}

func (s *SignatureService) GetSignatureByID(signatureID string) (common.Signature, error) {
	signatureDTO, err := s.repo.GetSignatureByID(signatureID)
	if errors.Is(err, persistence.ErrNotFound) {
		return common.Signature{}, ErrSignatureNotFound
	}

	if err != nil {
		return common.Signature{}, err
	}

	return signatureDTO.ToSignature(), nil
}

func (s *SignatureService) ListSignatures() ([]common.Signature, error) {
	signaturesDTOs, err := s.repo.ListSignatures()
	if err != nil {
		return nil, err
	}

	result := make([]common.Signature, len(signaturesDTOs))
	for i, sigDTO := range signaturesDTOs {
		result[i] = sigDTO.ToSignature()
	}
	return result, nil
}
