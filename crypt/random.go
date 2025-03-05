package crypt

import (
	"math/rand"
	"time"
)

type (
	CharSet string
)

const (
	Alphabet     CharSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	AlphaNum     CharSet = Alphabet + "0123456789"
	AlphaSpecial CharSet = AlphaNum + "~`!@#$%^&*()_-+={[}]|\\:;\"'<,>.?/"
)

var (
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func MakeRandom(length int, charset CharSet, prefix string) string {
	if length <= 0 {
		return ""
	}

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return prefix + string(b)
}
