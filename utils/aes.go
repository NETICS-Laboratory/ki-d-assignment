package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

func pad(input []byte, blockSize int) []byte {
	r := len(input) % blockSize
	pl := blockSize - r
	for i := 0; i < pl; i++ {
		input = append(input, byte(pl))
	}
	return input
}

func unpad(input []byte) ([]byte, error) {
	if len(input) == 0 {
		return nil, nil
	}

	pc := input[len(input)-1]
	pl := int(pc)

	if pl > len(input) || pl == 0 {
		return nil, errors.New("invalid padding")
	}

	p := input[len(input)-(pl):]
	for _, pc := range p {
		if uint(pc) != uint(len(p)) {
			return nil, errors.New("invalid padding")
		}
	}

	return input[:len(input)-pl], nil
}

func aesCbcEncrypt(plaintext, key []byte) ([]byte, error) {
	// Check if the plaintext size is not a multiple of the AES block size (16 bytes)
	if len(plaintext)%aes.BlockSize != 0 {
		return nil, errors.New("plaintext is not a multiple of the block size")
	}

	// Create an instance of the AES cipher using the provided key
	// The key size determines the security level (AES-128, AES-192, AES-256)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Prepare the ciphertext by allocating a byte slice to store both the IV and the encrypted data
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	// Generate a random IV using a cryptographically secure random number generator
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// Encrypt the plaintext using CBC mode (Cipher Block Chaining).
	// CBC mode requires the IV for chaining blocks of ciphertext.
	// The actual ciphertext will start after the IV in the ciphertext slice (i.e., at aes.BlockSize).
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func aesCbcDecrypt(ciphertext, key []byte) ([]byte, error) {
	// Create an instance of the AES cipher using the provided key
	// The key size determines the security level (AES-128, AES-192, AES-256)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	// Extract the Initialization Vector (IV) from the first 16 bytes of the ciphertext.
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// Check if the remaining ciphertext (after the IV) is a multiple of the AES block size.
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	// Initialize the AES CBC decryption mode using the block cipher and the IV.
	mode := cipher.NewCBCDecrypter(block, iv)

	// Perform the decryption of the ciphertext into the plaintext.
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	return plaintext, nil
}

func AESEncrypt(plaintext string, key []byte) (string, error) {
	plaintextbyte := []byte(plaintext)
	pad_plaintext := pad(plaintextbyte, 16)

	ciphertext, err := aesCbcEncrypt(pad_plaintext, key)
	if err != nil {
		return "", err
	}

	// Encode ciphertext to Base64 string
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func AESDecrypt(ciphertext string, key []byte) (string, error) {
	ciphertextbyte, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	pad_decrypt, err := aesCbcDecrypt(ciphertextbyte, key)
	if err != nil {
		return "", err
	}

	decrypt, err := unpad(pad_decrypt)
	if err != nil {
		return "", err
	}

	return string(decrypt), nil
}

func EncryptFileBytesAES(fileData []byte, key []byte) ([]byte, error) {
	padFileData := pad(fileData, aes.BlockSize)

	encryptedData, err := aesCbcEncrypt(padFileData, key)
	if err != nil {
		return nil, err
	}

	return encryptedData, nil
}

func DecryptFileBytesAES(encryptedFileBytes []byte, key []byte) ([]byte, error) {
	decryptedData, err := aesCbcDecrypt(encryptedFileBytes, key)
	if err != nil {
		return nil, err
	}

	unpadFileData, err := unpad(decryptedData)
	if err != nil {
		return nil, err
	}

	return unpadFileData, nil
}
