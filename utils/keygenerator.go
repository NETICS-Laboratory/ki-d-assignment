package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func GenerateAsymmetricKeys(username string) error {
	// Generate RSA key pair
	privateKey, publicKey, err := GenerateKeyPair()
	if err != nil {
		return err
	}

	// Define directories based on username
	privateKeyDirectory := fmt.Sprintf("uploads/%s/secret", username)
	publicKeyDirectory := fmt.Sprintf("uploads/%s/secret", username)

	// Create directories if they do not exist
	if err := os.MkdirAll(privateKeyDirectory, os.ModePerm); err != nil {
		return err
	}

	// Encode the private key
	privateKeyDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyDER,
	}
	privateKeyPath := filepath.Join(privateKeyDirectory, "private_key.pem")
	if err := os.WriteFile(privateKeyPath, pem.EncodeToMemory(privateKeyBlock), 0600); err != nil {
		return err
	}

	// Encode the public key
	publicKeyDER := x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyDER,
	}
	publicKeyPath := filepath.Join(publicKeyDirectory, "public_key.pem")
	if err := os.WriteFile(publicKeyPath, pem.EncodeToMemory(publicKeyBlock), 0644); err != nil {
		return err
	}

	return nil
}
