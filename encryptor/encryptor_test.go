package main

import (
	"bytes"
	"os"
	"testing"
)

func TestEncryptAndDecrypt(t *testing.T) {
	// Prepare
	input := []byte("This is a test WAV-like data.")
	inputFile := "test_input.wav"
	encFile := "test_output.enc"
	decFile := "test_output_decrypted.wav"

	err := os.WriteFile(inputFile, input, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(inputFile)

	// Encrypt
	err = EncryptFile(inputFile, encFile)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}
	defer os.Remove(encFile)

	// Decrypt
	err = DecryptFile(encFile, decFile)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}
	defer os.Remove(decFile)

	// Check result
	output, err := os.ReadFile(decFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(input, output) {
		t.Errorf("Decrypted content mismatch. Got %q, want %q", output, input)
	}
}
