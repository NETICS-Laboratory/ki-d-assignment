package utils

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"errors"
)

// Padding for DES block size (8 bytes)
func padDES(src []byte) []byte {
	padding := des.BlockSize - len(src)%des.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// Unpadding the decrypted data
func unpadDES(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

// DES Encryption function
func desEncryptutils(plaintext, key []byte) (string, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintext = padDES(plaintext)
	ciphertext := make([]byte, len(plaintext))

	iv := key[:des.BlockSize] // DES uses 8 bytes block size, use part of the key as IV

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	// Encode ciphertext to base64 for readable output
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DES Decryption function
func desDecryptutils(ciphertextBase64 string, key []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}

	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < des.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := key[:des.BlockSize]
	mode := cipher.NewCBCDecrypter(block, iv)

	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	plaintext = unpadDES(plaintext)

	return string(plaintext), nil
}

func DESEncrypt(plaintext string, key []byte) (string, error) {
	plaintextBytes := []byte(plaintext)

	ciphertext, err := desEncryptutils(plaintextBytes, key)
	if err != nil {
		return "", err
	}

	return ciphertext, nil
}

func DESDecrypt(ciphertext string, key []byte) (string, error) {
	plaintext, err := desDecryptutils(ciphertext, key)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func EncryptFileBytesDES(fileData []byte, key []byte) ([]byte, error) {
	padFileData := padDES(fileData)

	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := key[:des.BlockSize]
	ciphertext := make([]byte, len(padFileData))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padFileData)

	return ciphertext, nil
}

func DecryptFileBytesDES(encryptedFileBytes []byte, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(encryptedFileBytes) < des.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := key[:des.BlockSize]
	plaintext := make([]byte, len(encryptedFileBytes))

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, encryptedFileBytes)

	plaintext = unpadDES(plaintext)

	return plaintext, nil
}
