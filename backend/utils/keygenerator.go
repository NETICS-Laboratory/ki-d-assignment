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

func GetRSAPublicKey(user string) (*rsa.PublicKey, error) {
	// Get the user's public key
	// fmt.Println("Getting public key for user:", user)
	publicKeyPath := fmt.Sprintf("uploads/%s/secret/public_key.pem", user)

	publicKeyFile, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	publicKeyBlock, _ := pem.Decode(publicKeyFile)
	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

func PublicKeyToPEMString(pubKey *rsa.PublicKey) (string, error) {
	// Convert the public key to DER-encoded PKIX format
	pubKeyDER, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key: %v", err)
	}

	// Create a PEM block for the public key
	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyDER,
	})

	// Convert PEM block to a string
	return string(pubKeyPEM), nil
}
