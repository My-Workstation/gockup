package utils

import (
	"crypto/rand"
	"io"
)

func MakeRandom(length int) []byte {
	rnd := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, rnd); err != nil {
		panic(err)
	}
	return rnd
}
