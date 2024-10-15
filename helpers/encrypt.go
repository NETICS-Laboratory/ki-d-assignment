package helpers

import (
	"fmt"
	"ki-d-assignment/utils"
	"time"
)

func EncryptData(data string, secretKey []byte, secretKey8Byte []byte) (string, string, string, error) {
	startAES := time.Now()
	encAES, err := utils.AESEncrypt(data, secretKey)
	if err != nil {
		return "", "", "", err
	}
	fmt.Printf("AES Encryption Time: %.6f ms\n", float64(time.Since(startAES).Microseconds())/1000)
	fmt.Printf("AES Ciphertext Length: %d\n\n", len(encAES))

	startDES := time.Now()
	encDES, err := utils.DESEncrypt(data, secretKey8Byte)
	if err != nil {
		return "", "", "", err
	}
	fmt.Printf("DES Encryption Time: %.6f ms\n", float64(time.Since(startDES).Microseconds())/1000)
	fmt.Printf("DES Ciphertext Length: %d\n\n", len(encDES))

	startRC4 := time.Now()
	encRC4, err := utils.RC4Encrypt(data, secretKey)
	if err != nil {
		return "", "", "", err
	}
	fmt.Printf("RC4 Encryption Time: %.6f ms\n", float64(time.Since(startRC4).Microseconds())/1000)
	fmt.Printf("RC4 Ciphertext Length: %d\n\n", len(encRC4))

	return encAES, encDES, encRC4, nil
}
