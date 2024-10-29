package utils

import (
	"crypto/rand"
	"io"
	"crypto/rsa"
	"crypto/x509"
	"os"
	"encoding/pem"

	"github.com/google/uuid"
)

func GenerateSecretKey() ([]byte, error) {
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

func GenerateSecretKey8Byte() ([]byte, error) {
	key := make([]byte, 8) //
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func GenerateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

func GenerateAsymmetricKeys(id uuid.UUID) (error) {
	privateKey, publicKey, err := GenerateKeyPair()
	if err != nil {
		return err
	}
	
	privateKeyDirectory := "keys/private_keys"
	publicKeyDirectory := "keys/public_keys"

	if err := os.MkdirAll(privateKeyDirectory, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(publicKeyDirectory, os.ModePerm); err != nil {
		return err
	}

	privateKeyDER := x509.MarshalPKCS1PrivateKey(privateKey)
	publicKeyDER := x509.MarshalPKCS1PublicKey(publicKey)

	privateKeyBlock := &pem.Block{
		Type: "RSA PRIVATE KEY",
		Bytes: privateKeyDER,
	}
	
	publicKeyBlock := &pem.Block{
		Type: "RSA PUBLIC KEY",
		Bytes: publicKeyDER,
	}

	privateKeyName := privateKeyDirectory + "/" + id.String() + ".pem"
	publicKeyName := publicKeyDirectory + "/" + id.String() + ".pem"

	if err:= os.WriteFile(privateKeyName, pem.EncodeToMemory(privateKeyBlock), 0644); err != nil {
		return err
	}
	if err:= os.WriteFile(publicKeyName, pem.EncodeToMemory(publicKeyBlock), 0644); err != nil {
		return err
	}

	return nil
}