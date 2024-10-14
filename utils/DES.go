package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"fmt"
)

// Padding for DES block size (8 bytes)
func pad(src []byte) []byte {
	padding := des.BlockSize - len(src)%des.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// Unpadding the decrypted data
func unpad(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

// DES Encryption function
func desEncrypt(plaintext, key []byte) (string, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintext = pad(plaintext)
	ciphertext := make([]byte, len(plaintext))

	iv := key[:des.BlockSize] // DES uses 8 bytes block size, use part of the key as IV

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	// Encode ciphertext to base64 for readable output
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DES Decryption function
func desDecrypt(ciphertextBase64 string, key []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}

	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < des.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := key[:des.BlockSize]
	mode := cipher.NewCBCDecrypter(block, iv)

	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	plaintext = unpad(plaintext)

	return string(plaintext), nil
}

func main() {
	key := []byte("12345678") // DES key must be 8 bytes
	plaintext := "Hello, World!"

	// Encrypt
	encrypted, err := desEncrypt([]byte(plaintext), key)
	if err != nil {
		fmt.Println("Error encrypting:", err)
		return
	}
	fmt.Println("Encrypted (Base64):", encrypted)

	// Decrypt
	decrypted, err := desDecrypt(encrypted, key)
	if err != nil {
		fmt.Println("Error decrypting:", err)
		return
	}
	fmt.Println("Decrypted:", decrypted)
}