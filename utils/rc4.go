package utils

import (
	"crypto/rc4"
	"encoding/base64"
)

func rc4Encryptutils(plaintext, key []byte) ([]byte, error) {
	// create instance of rc4 cipher
	c, err := rc4.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// initialize ciphertext with same size
	ciphertext := make([]byte, len(plaintext))

	// execute XORKeyStream to encrypt plaintext
	c.XORKeyStream(ciphertext, plaintext)

	return ciphertext, nil
}

func rc4Decryptutils(ciphertext, key []byte) ([]byte, error) {
	// create instance of rc4 cipher
	c, err := rc4.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// initialize plaintext with same size
	plaintext := make([]byte, len(ciphertext))

	// execute XORKeyStream to decrypt ciphertext
	c.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

func RC4Encrypt(plaintext string, key []byte) (string, error) {
	plaintextBytes := []byte(plaintext)

	ciphertext, err := rc4Encryptutils(plaintextBytes, key)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func RC4Decrypt(ciphertext string, key []byte) (string, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	plaintext, err := rc4Decryptutils(ciphertextBytes, key)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func EncryptFileBytesRC4(fileData []byte, key []byte) ([]byte, error) {
	encryptedData, err := rc4Encryptutils(fileData, key)
	if err != nil {
		return nil, err
	}

	return encryptedData, nil
}

func DecryptFileBytesRC4(encryptedFileBytes []byte, key []byte) ([]byte, error) {
	decryptedData, err := rc4Decryptutils(encryptedFileBytes, key)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}
