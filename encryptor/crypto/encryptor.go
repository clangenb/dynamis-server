package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"os"
)

// LoadMasterKey loads the encryption key from the environment variable
func LoadMasterKey() ([]byte, error) {
	key := os.Getenv("AUDIO_ENCRYPTION_KEY")
	if len(key) != 32 {
		return nil, errors.New("master key must be 32 bytes long")
	}
	return []byte(key), nil
}

// EncryptFile encrypts an input file and writes to output using AES-256-GCM (streaming)
func EncryptFile(inputPath, outputPath string, masterKey []byte) error {
	plaintext, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer plaintext.Close()

	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	// Create output file and write the nonce as the first part
	ciphertextFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer ciphertextFile.Close()

	if _, err := ciphertextFile.Write(nonce); err != nil {
		return err
	}

	// Encrypt the file in chunks
	buf := make([]byte, 4096) // buffer size
	for {
		n, err := plaintext.Read(buf)
		if n > 0 {
			enc := gcm.Seal(nil, nonce, buf[:n], nil)
			if _, err := ciphertextFile.Write(enc); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// DecryptFile decrypts an encrypted file and writes to output using AES-256-GCM (streaming)
func DecryptFile(encPath, outputPath string, masterKey []byte) error {
	ciphertext, err := os.Open(encPath)
	if err != nil {
		return err
	}
	defer ciphertext.Close()

	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := ciphertext.Read(nonce); err != nil {
		return err
	}

	plaintextFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer plaintextFile.Close()

	buf := make([]byte, 4096) // buffer size
	for {
		n, err := ciphertext.Read(buf)
		if n > 0 {
			plaintext, err := gcm.Open(nil, nonce, buf[:n], nil)
			if err != nil {
				return err
			}
			if _, err := plaintextFile.Write(plaintext); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}
