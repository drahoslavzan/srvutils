package crypt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type HMACSignerParser struct {
	key []byte
	exp time.Duration
}

func NewHMACSignerParser(symmetricKey []byte, exp time.Duration) *HMACSignerParser {
	return &HMACSignerParser{
		key: symmetricKey,
		exp: exp,
	}
}

func (m *HMACSignerParser) Sign(payload JWTClaims) (string, error) {
	claims := makeClaims(payload, m.exp)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.key)
}

func (m *HMACSignerParser) Parse(token string) (JWTClaims, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.key, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return JWTClaims(claims), nil
	}

	return nil, errors.New("invalid token")
}
