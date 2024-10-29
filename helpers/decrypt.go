package helpers

import (
	"fmt"
	"ki-d-assignment/utils"
	"time"
)

func DecryptData(aesciphertext string, rc4ciphertext string, desciphertext string, secretKey []byte, secretKey8Byte []byte) (string, error) {
	const iterations = 1000
	var decAES, decDES, decRC4 string
	var err error

	startAES := time.Now()
	for i := 0; i < iterations; i++ {
		decAES, err = utils.AESDecrypt(aesciphertext, secretKey)
		if err != nil {
			return "", err
		}
	}
	fmt.Printf("AES Decryption Time: %.6f ms\n", float64(time.Since(startAES).Microseconds())/1000)
	fmt.Printf("AES Ciphertext: %v\n", aesciphertext)
	fmt.Printf("AES Ciphertext Length: %d\n\n", len(aesciphertext))

	startRC4 := time.Now()
	for i := 0; i < iterations; i++ {
		decRC4, err = utils.RC4Decrypt(rc4ciphertext, secretKey)
		if err != nil {
			return "", err
		}
	}
	fmt.Printf("RC4 Decryption Time: %.6f ms\n", float64(time.Since(startRC4).Microseconds())/1000)
	fmt.Printf("RC4 Ciphertext: %v\n", rc4ciphertext)
	fmt.Printf("RC4 Ciphertext Length: %d\n\n", len(rc4ciphertext))

	startDES := time.Now()
	for i := 0; i < iterations; i++ {
		decDES, err = utils.DESDecrypt(desciphertext, secretKey8Byte)
		if err != nil {
			return "", err
		}
	}
	fmt.Printf("DES Decryption Time: %.6f ms\n\n", float64(time.Since(startDES).Microseconds())/1000)
	fmt.Printf("DES Ciphertext: %v\n", desciphertext)
	fmt.Printf("DES Ciphertext Length: %d\n\n", len(desciphertext))

	if decAES == decDES && decDES == decRC4 {
		return decAES, nil
	}

	return "", err
}

func DecryptDataReturnIndiviual(aesciphertext string, rc4ciphertext string, desciphertext string, secretKey []byte, secretKey8Byte []byte) (string, string, string, error) {
	const iterations = 1000
	var decAES, decDES, decRC4 string
	var err error

	startAES := time.Now()
	for i := 0; i < iterations; i++ {
		decAES, err = utils.AESDecrypt(aesciphertext, secretKey)
		if err != nil {
			return "", "", "", err
		}
	}
	fmt.Printf("AES Decryption Time: %.6f ms\n", float64(time.Since(startAES).Microseconds())/1000)
	fmt.Printf("AES Ciphertext: %v\n", aesciphertext)
	fmt.Printf("AES Ciphertext Length: %d\n\n", len(aesciphertext))

	startRC4 := time.Now()
	for i := 0; i < iterations; i++ {
		decRC4, err = utils.RC4Decrypt(rc4ciphertext, secretKey)
		if err != nil {
			return "", "", "", err
		}
	}
	fmt.Printf("RC4 Decryption Time: %.6f ms\n", float64(time.Since(startRC4).Microseconds())/1000)
	fmt.Printf("RC4 Ciphertext: %v\n", rc4ciphertext)
	fmt.Printf("RC4 Ciphertext Length: %d\n\n", len(rc4ciphertext))

	startDES := time.Now()
	for i := 0; i < iterations; i++ {
		decDES, err = utils.DESDecrypt(desciphertext, secretKey8Byte)
		if err != nil {
			return "", "", "", err
		}
	}
	fmt.Printf("DES Decryption Time: %.6f ms\n\n", float64(time.Since(startDES).Microseconds())/1000)
	fmt.Printf("DES Ciphertext: %v\n", desciphertext)
	fmt.Printf("DES Ciphertext Length: %d\n\n", len(desciphertext))

	return decAES, decRC4, decDES, nil
}
