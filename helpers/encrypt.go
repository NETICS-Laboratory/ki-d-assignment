package helpers

import (
	"ki-d-assignment/utils"
)

func EncryptData(data string, secretKey []byte, secretKey8Byte []byte) (string, string, string, error) {
	encAES, err := utils.AESEncrypt(data, secretKey)
	if err != nil {
		return "", "", "", err
	}

	encDES, err := utils.DESEncrypt(data, secretKey8Byte)
	if err != nil {
		return "", "", "", err
	}

	encRC4, err := utils.RC4Encrypt(data, secretKey)
	if err != nil {
		return "", "", "", err
	}

	return encAES, encDES, encRC4, nil
}
