package utils

import (
	"crypto/rand"
	"io"
)

func GeneraretSecretKey() ([]byte, error) {
	// initiate secret key
	// key := make([]byte, 16) // AES-128
	// key := make([]byte, 24) // AES-192
	key := make([]byte, 32) // for AES-256 and RC4 later
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func GeneraretSecretKey8Byte() ([]byte, error) {
	key := make([]byte, 8) //
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
