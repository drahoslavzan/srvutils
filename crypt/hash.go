package crypt

import (
	"golang.org/x/crypto/bcrypt"
)

type Hash string

func MakeHash(text string) Hash {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return Hash(hash)
}

func (m Hash) Equals(plain string) bool {
	byteHash := []byte(m)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plain))

	return err == nil
}
