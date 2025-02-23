package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

// Encrypt data using AES-GCM
func EncryptAES(text string, key []byte) (ct string, nonce string, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}

	nb := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nb); err != nil {
		return "", "", err
	}

	cipherText := aesGCM.Seal(nil, nb, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(cipherText), base64.StdEncoding.EncodeToString(nb), nil
}

// Decrypt data using AES-GCM
func DecryptAES(ct, nonce string, key []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ct)
	if err != nil {
		return "", err
	}

	nonceBytes, err := base64.StdEncoding.DecodeString(nonce)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesGCM.Open(nil, nonceBytes, data, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
