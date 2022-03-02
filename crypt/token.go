package crypt

import (
	"github.com/google/uuid"
)

type Token string

func GenerateToken() Token {
	guid := uuid.New()
	return Token(guid.String())
}
