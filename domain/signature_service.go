package domain

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/common"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
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
	signature, err := s.repo.GetSignatureByID(signatureID)
	if err != nil {
		return common.Signature{}, err
	}

	return common.Signature{
		ID:         signature.ID,
		DeviceID:   signature.DeviceID,
		Signature:  signature.Signature,
		SignedData: signature.SignedData,
	}, nil
}

func (s *SignatureService) ListSignatures() ([]common.Signature, error) {
	signatures, err := s.repo.ListSignatures()
	if err != nil {
		return nil, err
	}

	result := make([]common.Signature, len(signatures))
	for i, sig := range signatures {
		result[i] = common.Signature{
			ID:         sig.ID,
			DeviceID:   sig.DeviceID,
			Signature:  sig.Signature,
			SignedData: sig.SignedData,
		}
	}
	return result, nil
}
