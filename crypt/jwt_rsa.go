package crypt

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
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

func NewRSASigner(privateKey []byte, exp time.Duration) *RSASigner {
	priKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		zap.L().Panic("invalid private key", zap.Error(err))
	}

	return &RSASigner{
		priKey: priKey,
	}
}

func (m *RSASigner) Sign(payload JWTClaims) (string, error) {
	claims := makeClaims(payload, m.exp)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(m.priKey)
}

func NewRSAParser(publicKey []byte) *RSAParser {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		zap.L().Panic("invalid public key", zap.Error(err))
	}

	return &RSAParser{
		pubKey: pubKey,
	}
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
