package persistence

import (
	"sync"

	"github.com/AloveIs/signing-device-service-go/common"
)

// InMemorySignatureDb implements an in-memory database for storing signature records
// using a map with read-write mutex for concurrent access control
type InMemorySignatureDb struct {
	// RWMutex to emulate atomicity of the database
	rwmutex sync.RWMutex
	// Storage method is a map signatureID:signature
	db map[string]common.SignatureDTO
}

func (db *InMemorySignatureDb) SaveSignature(signature common.SignatureDTO) error {
	db.rwmutex.Lock()
	defer db.rwmutex.Unlock()

	db.db[signature.ID] = signature
	return nil
}

func (db *InMemorySignatureDb) GetSignatureByID(signatureID string) (common.SignatureDTO, error) {
	db.rwmutex.RLock()
	defer db.rwmutex.RUnlock()

	if signature, ok := db.db[signatureID]; ok {
		return signature, nil
	}
	return common.SignatureDTO{}, ErrNotFound
}

func (db *InMemorySignatureDb) ListSignatures() ([]common.SignatureDTO, error) {
	db.rwmutex.RLock()
	defer db.rwmutex.RUnlock()

	signatures := make([]common.SignatureDTO, 0, len(db.db))
	for _, signature := range db.db {
		signatures = append(signatures, signature)
	}
	return signatures, nil
}

func (db *InMemorySignatureDb) GetSignaturesByDeviceID(deviceID string) ([]common.SignatureDTO, error) {
	db.rwmutex.RLock()
	defer db.rwmutex.RUnlock()

	var signatures []common.SignatureDTO
	for _, signature := range db.db {
		if signature.DeviceID == deviceID {
			signatures = append(signatures, signature)
		}
	}
	return signatures, nil
}

func NewInMemorySignatureDb() SignatureRepository {
	return &InMemorySignatureDb{
		db: make(map[string]common.SignatureDTO),
	}
}
