package server

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
)

// Encryptor handles all encryption-related operations.
type Encryptor struct {
	MasterKey []byte
}

type KeyDeriver interface {
	DeriveKey(audioID, userID, deviceID string) ([]byte, error)
}

// NewEncryptor initializes a new Encryptor instance with a master key.
func NewEncryptor(masterKey []byte) *Encryptor {
	return &Encryptor{MasterKey: masterKey}
}

// LoadMasterKey loads the master key from an environment variable.
func (e *Encryptor) LoadMasterKey() ([]byte, error) {
	if len(e.MasterKey) == 0 {
		return nil, errors.New("master key is not set")
	}
	return e.MasterKey, nil
}

// DeriveKey generates a derived key based on the audio ID, user ID, and device ID.
func (e *Encryptor) DeriveKey(audioID, userID, deviceID string) ([]byte, error) {
	// Use SHA-256 as a deterministic and fixed-length seed for key derivation.
	input := []byte(audioID + userID + deviceID)
	hash := sha256.Sum256(input)

	block, err := aes.NewCipher(e.MasterKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := hash[:gcm.NonceSize()] // Always safe (nonce size is 12)
	return gcm.Seal(nil, nonce, hash[:], nil)[:32], nil
}
