package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func GenerateEncryptedHash (data string, privateKey *rsa.PrivateKey) (string, error) {
	hashed := sha256.Sum256([]byte(data))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(signature), nil
}

func VerifySignature (data string, signature string, publicKey *rsa.PublicKey) error {
	hashed := sha256.Sum256([]byte(data))
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return err
	}
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signatureBytes)
	if err != nil {
		return err
	}
	return nil
}

func ParsePublicKey