package crypt

import (
	"strings"

	"github.com/google/uuid"
)

type Token string

func GenerateToken() Token {
	guid := uuid.New()
	return Token(guid.String())
}

func (m Token) ToID() string {
	return strings.ReplaceAll(string(m), "-", "")
}
