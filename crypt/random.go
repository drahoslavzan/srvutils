package crypt

import (
	"math/rand"
	"time"
)

var (
	Alphabet     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	AlphaNum     = Alphabet + "0123456789"
	AlphaSpecial = AlphaNum + "~`!@#$%^&*()_-+={[}]|\\:;\"'<,>.?/"

	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func MakeRandom(length int, charset string, prefix string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset)-1)]
	}

	return prefix + string(b)
}
