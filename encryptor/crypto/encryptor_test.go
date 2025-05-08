package crypto

import (
	"bytes"
	"os"
	"strconv"
	"testing"
)

func testMasterKey() []byte {
	key := []byte("TestMasterKey-------------------")
	if len(key) != 32 {
		panic("invalid test key len" + strconv.Itoa(len(key)))
	}

	return key
}

func TestEncryptAndDecrypt(t *testing.T) {
	// Prepare
	input := []byte("This is some audio-like data.")
	inputFile := "test_input.wav"
	encFile := "test_output.enc"
	decFile := "test_output_decrypted.wav"

	err := os.WriteFile(inputFile, input, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(inputFile)

	// Get the master key from the environment
	masterKey := testMasterKey()

	// Encrypt the file
	err = EncryptFile(inputFile, encFile, masterKey)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}
	defer os.Remove(encFile)

	// Decrypt the file
	err = DecryptFile(encFile, decFile, masterKey)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}
	defer os.Remove(decFile)

	// Check if decrypted file content matches the original
	decryptedContent, err := os.ReadFile(decFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(input, decryptedContent) {
		t.Errorf("Decrypted content mismatch. Got %q, want %q", decryptedContent, input)
	}
}
