package crypt

import (
	"golang.org/x/crypto/bcrypt"
)

type Hash []byte

func MakeHash(text string) (Hash, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return Hash(hash), nil
}

func (m Hash) Equals(plain string) bool {
	return bcrypt.CompareHashAndPassword(m, []byte(plain)) == nil
}
