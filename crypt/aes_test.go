package crypt

import (
	"crypto/rand"
	"testing"
)

func TestEncryptDecryptString(t *testing.T) {
	// Test key (32 bytes for AES-256)
	key := []byte("32-byte-long-key-134567890123456")

	// Test plaintext
	plaintext := "This is a secret message."

	// Encrypt the plaintext
	ciphertext, err := EncryptAES(plaintext, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Decrypt the ciphertext
	decrypted, err := DecryptAES(ciphertext, key)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	// Check if the decrypted text matches the original plaintext
	if decrypted != plaintext {
		t.Errorf("Decrypted text does not match original plaintext. Got: %s, Want: %s", decrypted, plaintext)
	}
}

func TestEncryptString_InvalidKey(t *testing.T) {
	// Invalid key (not 16, 24, or 32 bytes)
	invalidKey := []byte("short-key")

	// Test plaintext
	plaintext := "This is a secret message."

	// Attempt to encrypt with an invalid key
	_, err := EncryptAES(plaintext, invalidKey)
	if err == nil {
		t.Error("Expected encryption to fail with invalid key, but it succeeded")
	}
}

func TestDecryptString_InvalidCiphertext(t *testing.T) {
	// Test key (32 bytes for AES-256)
	key := []byte("32-byte-long-key-134567890123456")

	// Invalid ciphertext (too short to contain a nonce)
	invalidCiphertext := []byte("short")

	// Attempt to decrypt invalid ciphertext
	_, err := DecryptAES(invalidCiphertext, key)
	if err == nil {
		t.Error("Expected decryption to fail with invalid ciphertext, but it succeeded")
	}
}

func TestDecryptString_TamperedCiphertext(t *testing.T) {
	// Test key (32 bytes for AES-256)
	key := []byte("32-byte-long-key-123456790123456")

	// Test plaintext
	plaintext := "This is a secret message."

	// Encrypt the plaintext
	ciphertext, err := EncryptAES(plaintext, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Tamper with the ciphertext (e.g., modify a byte)
	tamperedCiphertext := string(ciphertext[:len(ciphertext)-1]) + "x"

	// Attempt to decrypt the tampered ciphertext
	_, err = DecryptAES([]byte(tamperedCiphertext), key)
	if err == nil {
		t.Error("Expected decryption to fail with tampered ciphertext, but it succeeded")
	}
}

func TestEncryptString_EmptyPlaintext(t *testing.T) {
	// Test key (32 bytes for AES-256)
	key := []byte("32-byte-lng-key-1234567890123456")

	// Empty plaintext
	plaintext := ""

	// Encrypt the empty plaintext
	ciphertext, err := EncryptAES(plaintext, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Decrypt the ciphertext
	decrypted, err := DecryptAES(ciphertext, key)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	// Check if the decrypted text matches the original plaintext
	if decrypted != plaintext {
		t.Errorf("Decrypted text does not match original plaintext. Got: %s, Want: %s", decrypted, plaintext)
	}
}

func TestEncryptString_RandomKey(t *testing.T) {
	// Generate a random key (32 bytes for AES-256)
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		t.Fatalf("Failed to generate random key: %v", err)
	}

	// Test plaintext
	plaintext := "This is a secret message."

	// Encrypt the plaintext
	ciphertext, err := EncryptAES(plaintext, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Decrypt the ciphertext
	decrypted, err := DecryptAES(ciphertext, key)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	// Check if the decrypted text matches the original plaintext
	if decrypted != plaintext {
		t.Errorf("Decrypted text does not match original plaintext. Got: %s, Want: %s", decrypted, plaintext)
	}
}
