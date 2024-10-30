package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
)

func GetPublicKey(id uuid.UUID) (string, error) {
	publicKeyDirectory := "keys/public_keys"
	publicKeyPath := fmt.Sprintf("%s/%s.pem", publicKeyDirectory, id.String())
	publicKeyFile, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(publicKeyFile)
	if block == nil {
		return "", err
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(publicKey.N.Bytes()), nil
}