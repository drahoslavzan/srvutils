package crypt

import (
	"golang.org/x/crypto/bcrypt"
)

type Hash string

func MakeHash(str string) Hash {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return Hash(hash)
}

func (m Hash) IsEqualTo(plain string) bool {
	byteHash := []byte(m)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plain))
	if err != nil {
		return false
	}
	return true
}
