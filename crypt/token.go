package crypt

import (
	"github.com/google/uuid"
)

type Token string

func Generate() Token {
	guid := uuid.New()
	return Token(guid.String())
}
