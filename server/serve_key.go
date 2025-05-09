package server

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
	"os"
)

// Encryptor handles all encryption-related operations.
type Encryptor struct {
	MasterKey []byte
}

type KeyDeriver interface {
	DeriveCEK(audioID, userID, deviceID string) ([]byte, error)
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

// DeriveCEK generates a unique CEK using audioID, userID, deviceID.
func (e *Encryptor) DeriveCEK(audioID, userID, deviceID string) ([]byte, error) {
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

	nonce := hash[:gcm.NonceSize()]
	cek := gcm.Seal(nil, nonce, hash[:], nil)[:32] // 256-bit key
	return cek, nil
}

// EncryptFile encrypts a file with the CEK and stores it.
func (e *Encryptor) EncryptFile(inputPath, outputPath string, cek []byte) error {
	in, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(cek)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize()) // Securely generate this for real use
	copy(nonce, cek[:gcm.NonceSize()])     // For demo purposes only

	encrypted := gcm.Seal(nonce, nonce, in, nil)
	return os.WriteFile(outputPath, encrypted, 0644)
}
