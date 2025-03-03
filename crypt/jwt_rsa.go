package crypt

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (
	RSASigner struct {
		priKey *rsa.PrivateKey
		exp    time.Duration
	}

	RSAParser struct {
		pubKey *rsa.PublicKey
	}
)

func NewRSASigner(privateKey []byte, exp time.Duration) (*RSASigner, error) {
	priKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return nil, err
	}

	ret := &RSASigner{
		priKey: priKey,
		exp:    exp,
	}

	return ret, nil
}

func (m *RSASigner) Sign(payload JWTClaims) (string, error) {
	claims := makeClaims(payload, m.exp)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(m.priKey)
}

func NewRSAParser(publicKey []byte) (*RSAParser, error) {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, err
	}

	ret := &RSAParser{
		pubKey: pubKey,
	}

	return ret, nil
}

func (m *RSAParser) Parse(token string) (JWTClaims, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.pubKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return JWTClaims(claims), nil
	}

	return nil, errors.New("invalid token")
}
