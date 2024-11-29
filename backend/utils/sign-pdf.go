package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
)

func GenerateEncryptedHash(msg []byte, publicKey *rsa.PublicKey) (string, error) {
	msgHash := sha256.Sum256(msg)

	encryptedMsgHash, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, msgHash[:], nil)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encryptedMsgHash), nil
}

func VerifyDigitalSignature(signature string, fileSignature string) (bool, error) {

	return signature == fileSignature, nil
}
