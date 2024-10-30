package utils

import (
	"crypto/rand"
	"io"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"ki-d-assignment/dto"
)

func EncryptSymmetricKeyRSA(key []byte, publicKey *rsa.PublicKey) (string, error) {
	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, key, nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedKey), nil
}