package helpers

import (
	"ki-d-assignment/utils"
)

func DecryptData(aesciphertext string, rc4ciphertext string, desciphertext string, secretKey []byte, secretKey8Byte []byte) (string, error) {
	decAES, err := utils.AESDecrypt(aesciphertext, secretKey)
	if err != nil {
		return "", err
	}

	decRC4, err := utils.RC4Decrypt(rc4ciphertext, secretKey)
	if err != nil {
		return "", err
	}

	decDES, err := utils.DESDecrypt(desciphertext, secretKey8Byte)
	if err != nil {
		return "", err
	}

	if decAES == decDES && decDES == decRC4 {
		return decAES, nil
	}

	return "", err
}

func DecryptDataReturnIndiviual(aesciphertext string, rc4ciphertext string, desciphertext string, secretKey []byte, secretKey8Byte []byte) (string, string, string, error) {
	decAES, err := utils.AESDecrypt(aesciphertext, secretKey)
	if err != nil {
		return "", "", "", err
	}

	decRC4, err := utils.RC4Decrypt(rc4ciphertext, secretKey)
	if err != nil {
		return "", "", "", err
	}

	decDES, err := utils.DESDecrypt(desciphertext, secretKey8Byte)
	if err != nil {
		return "", "", "", err
	}

	return decAES, decRC4, decDES, nil
}
