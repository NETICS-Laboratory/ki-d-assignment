package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// IsValidFileName checks if the file name contains valid characters and extensions
func IsValidFileName(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext != ".jpeg" && ext != ".jpg" && ext != ".png" {
		return false
	}
	// Optionally, add more validations for file name, such as checking for illegal characters.
	return true
}

func removeExtension(fileName string, extension string) string {
	if filepath.Ext(fileName) == extension {
		return fileName[:len(fileName)-len(extension)]
	}
	return fileName
}

func UploadFile(fileHeader *multipart.FileHeader, filePath string, secretKey []byte, secretKey8Byte []byte) (string, string, string, error) {
	// Create user directory based on filePath (user.ID)
	if err := os.MkdirAll(filePath, 0755); err != nil {
		return "", "", "", err
	}

	// Open uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return "", "", "", err
	}
	defer file.Close()

	// Read file as byte array
	fileData, err := io.ReadAll(file)
	if err != nil {
		return "", "", "", err
	}

	// Encrypt file using AES, RC4, and DES
	encryptedFileAES, err := EncryptFileBytesAES(fileData, secretKey)
	if err != nil {
		return "", "", "", err
	}

	encryptedFileRC4, err := EncryptFileBytesRC4(fileData, secretKey)
	if err != nil {
		return "", "", "", err
	}

	encryptedFileDES, err := EncryptFileBytesDES(fileData, secretKey8Byte)
	if err != nil {
		return "", "", "", err
	}

	// Store encrypted AES file with .aes extension
	aesPath := filepath.Join(filePath, "aes")
	if err := os.MkdirAll(aesPath, 0755); err != nil {
		return "", "", "", err
	}
	aesFile := filepath.Join(aesPath, fileHeader.Filename+".aes")
	if err := os.WriteFile(aesFile, encryptedFileAES, 0644); err != nil {
		return "", "", "", err
	}

	// Store encrypted RC4 file with .rc4 extension
	rc4Path := filepath.Join(filePath, "rc4")
	if err := os.MkdirAll(rc4Path, 0755); err != nil {
		return "", "", "", err
	}
	rc4File := filepath.Join(rc4Path, fileHeader.Filename+".rc4")
	if err := os.WriteFile(rc4File, encryptedFileRC4, 0644); err != nil {
		return "", "", "", err
	}

	// Store encrypted DES file with .des extension
	desPath := filepath.Join(filePath, "des")
	if err := os.MkdirAll(desPath, 0755); err != nil {
		return "", "", "", err
	}
	desFile := filepath.Join(desPath, fileHeader.Filename+".des")
	if err := os.WriteFile(desFile, encryptedFileDES, 0644); err != nil {
		return "", "", "", err
	}

	// Return the paths of the encrypted files
	return aesFile, rc4File, desFile, nil
}

func DecryptAndSaveFiles(filePath string, aesFilePath string, rc4FilePath string, desFilePath string, secretKey []byte, secretKey8Byte []byte) error {
	// Create decrypted folder
	decryptedPath := filepath.Join(filePath, "decrypted")
	if err := os.MkdirAll(decryptedPath, 0755); err != nil {
		return err
	}

	// Decrypt and save AES file
	aesDecryptedFolder := filepath.Join(decryptedPath, "aes")
	if err := os.MkdirAll(aesDecryptedFolder, 0755); err != nil {
		return err
	}
	aesEncryptedData, err := os.ReadFile(aesFilePath)
	if err != nil {
		return err
	}
	decryptedAES, err := DecryptFileBytesAES(aesEncryptedData, secretKey)
	if err != nil {
		return err
	}
	decryptedAESPath := filepath.Join(aesDecryptedFolder, removeExtension(filepath.Base(aesFilePath), ".aes"))
	if err := os.WriteFile(decryptedAESPath, decryptedAES, 0644); err != nil {
		return err
	}

	// Decrypt and save RC4 file
	rc4DecryptedFolder := filepath.Join(decryptedPath, "rc4")
	if err := os.MkdirAll(rc4DecryptedFolder, 0755); err != nil {
		return err
	}
	rc4EncryptedData, err := os.ReadFile(rc4FilePath)
	if err != nil {
		return err
	}
	decryptedRC4, err := DecryptFileBytesRC4(rc4EncryptedData, secretKey)
	if err != nil {
		return err
	}
	decryptedRC4Path := filepath.Join(rc4DecryptedFolder, removeExtension(filepath.Base(rc4FilePath), ".rc4"))
	if err := os.WriteFile(decryptedRC4Path, decryptedRC4, 0644); err != nil {
		return err
	}

	// Decrypt and save DES file
	desDecryptedFolder := filepath.Join(decryptedPath, "des")
	if err := os.MkdirAll(desDecryptedFolder, 0755); err != nil {
		return err
	}
	desEncryptedData, err := os.ReadFile(desFilePath)
	if err != nil {
		return err
	}
	decryptedDES, err := DecryptFileBytesDES(desEncryptedData, secretKey8Byte)
	if err != nil {
		return err
	}
	decryptedDESPath := filepath.Join(desDecryptedFolder, removeExtension(filepath.Base(desFilePath), ".des"))
	if err := os.WriteFile(decryptedDESPath, decryptedDES, 0644); err != nil {
		return err
	}

	return nil
}

// Kalau mau uji decrypt tapi masih ada extentionnya
// func DecryptAndSaveFiles(filePath string, aesFilePath string, rc4FilePath string, desFilePath string, secretKey []byte, secretKey8Byte []byte) error {
// 	// Create decrypted folder
// 	decryptedPath := filepath.Join(filePath, "decrypted")
// 	if err := os.MkdirAll(decryptedPath, 0755); err != nil {
// 		return err
// 	}

// 	// Read and decrypt AES file
// 	aesEncryptedData, err := os.ReadFile(aesFilePath)
// 	if err != nil {
// 		return err
// 	}
// 	decryptedAES, err := DecryptFileBytesAES(aesEncryptedData, secretKey)
// 	if err != nil {
// 		return err
// 	}
// 	decryptedAESPath := filepath.Join(decryptedPath, "decrypted_aes_"+filepath.Base(aesFilePath))
// 	if err := os.WriteFile(decryptedAESPath, decryptedAES, 0644); err != nil {
// 		return err
// 	}

// 	// Read and decrypt RC4 file
// 	rc4EncryptedData, err := os.ReadFile(rc4FilePath)
// 	if err != nil {
// 		return err
// 	}
// 	decryptedRC4, err := DecryptFileBytesRC4(rc4EncryptedData, secretKey)
// 	if err != nil {
// 		return err
// 	}
// 	decryptedRC4Path := filepath.Join(decryptedPath, "decrypted_rc4_"+filepath.Base(rc4FilePath))
// 	if err := os.WriteFile(decryptedRC4Path, decryptedRC4, 0644); err != nil {
// 		return err
// 	}

// 	// Read and decrypt DES file
// 	desEncryptedData, err := os.ReadFile(desFilePath)
// 	if err != nil {
// 		return err
// 	}
// 	decryptedDES, err := DecryptFileBytesDES(desEncryptedData, secretKey8Byte)
// 	if err != nil {
// 		return err
// 	}
// 	decryptedDESPath := filepath.Join(decryptedPath, "decrypted_des_"+filepath.Base(desFilePath))
// 	if err := os.WriteFile(decryptedDESPath, decryptedDES, 0644); err != nil {
// 		return err
// 	}

// 	return nil
// }