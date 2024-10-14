package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"

	// "fmt"
	"io"
	"os"
	// "github.com/joho/godotenv"
)

// Padding plaintext to match block size (PKCS7 padding)
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// Remove padding after decryption (PKCS7 unpadding)
func pkcs7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	unpadding := int(data[length-1])
	if unpadding > length || unpadding == 0 {
		return nil, errors.New("invalid padding")
	}
	return data[:(length - unpadding)], nil
}

// Helper function to create a 32-byte AES key using a passphrase
func generateAESKey() []byte {
	key := sha256.Sum256([]byte(os.Getenv("KEY")))
	return key[:]
}

// AES CBC encryption function
func EncryptAESCBC(plaintext string) (string, error) {
	key := generateAESKey()

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintextBytes := pkcs7Padding([]byte(plaintext), aes.BlockSize)

	ciphertext := make([]byte, aes.BlockSize+len(plaintextBytes))
	iv := ciphertext[:aes.BlockSize] // Initialization Vector (IV)

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintextBytes)

	result := base64.URLEncoding.EncodeToString(ciphertext)

	return result, nil
}

func EncryptAESCBCFile(plainfile []byte) ([]byte, error) {
	key := generateAESKey()

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plainfileBytes := pkcs7Padding(plainfile, aes.BlockSize)

	cipherfile := make([]byte, aes.BlockSize+len(plainfileBytes))
	iv := cipherfile[:aes.BlockSize] // Initialization Vector (IV)

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherfile[aes.BlockSize:], plainfileBytes)

	return cipherfile, nil
}

// AES CBC decryption function
func DecryptAESCBC(ciphertext string) (string, error) {
	cipherdata, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	key := generateAESKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(cipherdata) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := cipherdata[:aes.BlockSize]
	cipherdata = cipherdata[aes.BlockSize:]

	if len(cipherdata)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherdata, cipherdata)

	plaintextBytes, err := pkcs7Unpadding(cipherdata)
	if err != nil {
		return "", err
	}

	result := string(plaintextBytes)

	return result, nil
}

func DecryptAESCBCFile(cipherfile []byte) ([]byte, error) {

	key := generateAESKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// if len(cipherdata) < aes.BlockSize {
	// 	return "", errors.New("cipherfile too short")
	// }

	iv := cipherfile[:aes.BlockSize]
	cipherfile = cipherfile[aes.BlockSize:]

	// if len(cipherdata)%aes.BlockSize != 0 {
	// 	return "", errors.New("cipherfile is not a multiple of the block size")
	// }

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherfile, cipherfile)

	plainfileBytes, err := pkcs7Unpadding(cipherfile)
	if err != nil {
		return nil, err
	}

	result := plainfileBytes

	return result, nil
}
