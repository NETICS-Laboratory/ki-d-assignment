package helpers

import (
	"fmt"
	"ki-d-assignment/utils"
	"time"
)

func DecryptData(aesciphertext string, rc4ciphertext string, desciphertext string, secretKey []byte, secretKey8Byte []byte) (string, error) {
	startAES := time.Now()
	decAES, err := utils.AESDecrypt(aesciphertext, secretKey)
	if err != nil {
		return "", err
	}
	fmt.Printf("AES Decryption Time: %.6f ms\n", float64(time.Since(startAES).Microseconds())/1000)

	startRC4 := time.Now()
	decRC4, err := utils.RC4Decrypt(rc4ciphertext, secretKey)
	if err != nil {
		return "", err
	}
	fmt.Printf("RC4 Decryption Time: %.6f ms\n", float64(time.Since(startRC4).Microseconds())/1000)

	startDES := time.Now()
	decDES, err := utils.DESDecrypt(desciphertext, secretKey8Byte)
	if err != nil {
		return "", err
	}
	fmt.Printf("DES Decryption Time: %.6f ms\n\n", float64(time.Since(startDES).Microseconds())/1000)

	if decAES == decDES && decDES == decRC4 {
		return decAES, nil
	}

	return "", err
}

func DecryptDataReturnIndiviual(aesciphertext string, rc4ciphertext string, desciphertext string, secretKey []byte, secretKey8Byte []byte) (string, string, string, error) {
	startAES := time.Now()
	decAES, err := utils.AESDecrypt(aesciphertext, secretKey)
	if err != nil {
		return "", "", "", err
	}
	fmt.Printf("AES Decryption Time: %.6f ms\n", float64(time.Since(startAES).Microseconds())/1000)

	startRC4 := time.Now()
	decRC4, err := utils.RC4Decrypt(rc4ciphertext, secretKey)
	if err != nil {
		return "", "", "", err
	}
	fmt.Printf("RC4 Decryption Time: %.6f ms\n", float64(time.Since(startRC4).Microseconds())/1000)

	startDES := time.Now()
	decDES, err := utils.DESDecrypt(desciphertext, secretKey8Byte)
	if err != nil {
		return "", "", "", err
	}
	fmt.Printf("DES Decryption Time: %.6f ms\n\n", float64(time.Since(startDES).Microseconds())/1000)

	return decAES, decRC4, decDES, nil
}
