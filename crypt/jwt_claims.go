package crypt

import (
	"time"

	"maps"

	"github.com/golang-jwt/jwt/v5"
)

type (
	JWTClaims jwt.MapClaims
)

func NewBasicJWTPayload(sub, aud, iss string, data any) JWTClaims {
	claims := JWTClaims{
		"sub": sub,
		"aud": aud,
		"iss": iss,
	}

	if data != nil {
		claims["data"] = data
	}

	return claims
}

func (m JWTClaims) Subject() (string, error) {
	return jwt.MapClaims(m).GetSubject()
}

func (m JWTClaims) Audience() ([]string, error) {
	return jwt.MapClaims(m).GetAudience()
}

func makeClaims(payload JWTClaims, exp time.Duration) jwt.MapClaims {
	now := time.Now()
	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(exp).Unix(),
	}

	maps.Copy(claims, payload)

	return claims
}
