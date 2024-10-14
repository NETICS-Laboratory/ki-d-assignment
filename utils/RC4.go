package utils

import (
    "crypto/rc4"
	"encoding/base64"
)

func RC4Encrypt(plaintext []byte) (string, error) {
    key := []byte(GetEnv("KEY"))
    
    cipher, err := rc4.NewCipher(key)
    if err != nil {
        return "", err 
    }

    ciphertext := make([]byte, len(plaintext))
    cipher.XORKeyStream(ciphertext, plaintext)

    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func RC4Decrypt(encrypted string) (string, error) {
    key := []byte(GetEnv("KEY"))
    
    cipher, err := rc4.NewCipher(key)
    if err != nil {
        return "", err 
    }

    encryptedBytes, err := base64.StdEncoding.DecodeString(encrypted)
    if err != nil {
        return "", err 
    }

    plaintext := make([]byte, len(encryptedBytes))
    cipher.XORKeyStream(plaintext, encryptedBytes)

    return string(plaintext), nil
}
