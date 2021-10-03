package crypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

type crypt struct {
	key *rsa.PrivateKey
}

type Crypt interface {
	SignPKCS1Base64(message string) string
}

func NewCrypt(pemPrivKey []byte) Crypt {
	privPem, _ := pem.Decode(pemPrivKey)
	if privPem.Type != "RSA PRIVATE KEY" {
		panic("invalid private key (PEM format expected)")
	}

	pk, err := x509.ParsePKCS1PrivateKey(privPem.Bytes)
	if err != nil {
		panic(err)
	}

	return &crypt{
		key: pk,
	}
}

func (m *crypt) SignPKCS1Base64(message string) string {
	hash := sha256.New()
	data := []byte(message)
	hash.Write(data)
	d := hash.Sum(nil)
	ctext, err := rsa.SignPKCS1v15(rand.Reader, m.key, crypto.SHA256, d)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(ctext)
}
