package utils

import (
	"crypto/rand"
	"io"
	"log"
	"os"
)

func MakeRandom(length int) []byte {
	rnd := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, rnd); err != nil {
		panic(err)
	}
	return rnd
}

func ReadKeyFromFile(filePath string) []byte {
	// try to open file
	encryptionKeyFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Can not open key file. %v", err)
	}
	defer encryptionKeyFile.Close()
	tmp := make([]byte, 33)
	count, err := encryptionKeyFile.Read(tmp)
	if err != nil {
		log.Fatal("Can not read key file. %v", err)
	}
	if count < 32 {
		log.Fatal("Key should be 32 bytes")
	} else if count > 32 {
		log.Print("Key should be 32 bytes. Only the first 32 bytes will be used.")
	}
	return tmp[:32]
}
