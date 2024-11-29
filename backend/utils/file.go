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
	// Allowed file extensions
	ext := strings.ToLower(filepath.Ext(fileName))
	validExtensions := map[string]bool{
		".jpeg": true,
		".jpg":  true,
		".png":  true,
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".xls":  true,
		".xlsx": true,
		".mp4":  true,
	}

	// Check if file extension is valid
	if _, valid := validExtensions[ext]; !valid {
		return false
	}

	// Optionally, add more validation for illegal characters
	illegalChars := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}
	for _, char := range illegalChars {
		if strings.Contains(fileName, char) {
			return false
		}
	}

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
	// startAES := time.Now()
	encryptedFileAES, err := EncryptFileBytesAES(fileData, secretKey)
	if err != nil {
		return "", "", "", err
	}
	// fmt.Printf("AES File Encryption Time: %.6f ms\n", float64(time.Since(startAES).Nanoseconds())/1e6)
	// fmt.Printf("AES Encrypted File Size: %d bytes\n\n", len(encryptedFileAES))

	// startRC4 := time.Now()
	encryptedFileRC4, err := EncryptFileBytesRC4(fileData, secretKey)
	if err != nil {
		return "", "", "", err
	}
	// fmt.Printf("RC4 File Encryption Time: %.6f ms\n", float64(time.Since(startRC4).Nanoseconds())/1e6)
	// fmt.Printf("RC4 Encrypted File Size: %d bytes\n\n", len(encryptedFileRC4))

	// startDES := time.Now()
	encryptedFileDES, err := EncryptFileBytesDES(fileData, secretKey8Byte)
	if err != nil {
		return "", "", "", err
	}
	// fmt.Printf("DES File Encryption Time: %.6f ms\n", float64(time.Since(startDES).Nanoseconds())/1e6)
	// fmt.Printf("DES Encrypted File Size: %d bytes\n\n", len(encryptedFileDES))

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
	// startAES := time.Now()
	decryptedAES, err := DecryptFileBytesAES(aesEncryptedData, secretKey)
	if err != nil {
		return err
	}
	// fmt.Printf("AES File Decryption Time: %.6f ms\n", float64(time.Since(startAES).Nanoseconds())/1e6)
	// fmt.Printf("AES Decrypted File Size: %d bytes\n\n", len(decryptedAES))

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
	// startRC4 := time.Now()
	decryptedRC4, err := DecryptFileBytesRC4(rc4EncryptedData, secretKey)
	if err != nil {
		return err
	}
	// fmt.Printf("RC4 File Decryption Time: %.6f ms\n", float64(time.Since(startRC4).Nanoseconds())/1e6)
	// fmt.Printf("RC4 Decrypted File Size: %d bytes\n\n", len(decryptedRC4))

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
	// startDES := time.Now()
	decryptedDES, err := DecryptFileBytesDES(desEncryptedData, secretKey8Byte)
	if err != nil {
		return err
	}
	// fmt.Printf("DES File Decryption Time: %.6f ms\n\n", float64(time.Since(startDES).Nanoseconds())/1e6)
	// fmt.Printf("DES Decrypted File Size: %d bytes\n\n", len(decryptedDES))

	decryptedDESPath := filepath.Join(desDecryptedFolder, removeExtension(filepath.Base(desFilePath), ".des"))
	if err := os.WriteFile(decryptedDESPath, decryptedDES, 0644); err != nil {
		return err
	}

	return nil
}

func DecryptAndSaveFilesReturnData(filePath string, aesFilePath string, rc4FilePath string, desFilePath string, secretKey []byte, secretKey8Byte []byte) (string, string, string, error) {
	// Create decrypted folder
	decryptedPath := filepath.Join(filePath, "decrypted")
	if err := os.MkdirAll(decryptedPath, 0755); err != nil {
		return "", "", "", err
	}

	// Decrypt and save AES file
	aesDecryptedFolder := filepath.Join(decryptedPath, "aes")
	if err := os.MkdirAll(aesDecryptedFolder, 0755); err != nil {
		return "", "", "", err
	}
	aesEncryptedData, err := os.ReadFile(aesFilePath)
	if err != nil {
		return "", "", "", err
	}
	// startAES := time.Now()
	decryptedAES, err := DecryptFileBytesAES(aesEncryptedData, secretKey)
	if err != nil {
		return "", "", "", err
	}
	// fmt.Printf("AES File Decryption Time: %.6f ms\n", float64(time.Since(startAES).Nanoseconds())/1e6)
	// fmt.Printf("AES Decrypted File Size: %d bytes\n\n", len(decryptedAES))

	decryptedAESPath := filepath.Join(aesDecryptedFolder, removeExtension(filepath.Base(aesFilePath), ".aes"))
	if err := os.WriteFile(decryptedAESPath, decryptedAES, 0644); err != nil {
		return "", "", "", err
	}

	// Decrypt and save RC4 file
	rc4DecryptedFolder := filepath.Join(decryptedPath, "rc4")
	if err := os.MkdirAll(rc4DecryptedFolder, 0755); err != nil {
		return "", "", "", err
	}
	rc4EncryptedData, err := os.ReadFile(rc4FilePath)
	if err != nil {
		return "", "", "", err
	}
	// startRC4 := time.Now()
	decryptedRC4, err := DecryptFileBytesRC4(rc4EncryptedData, secretKey)
	if err != nil {
		return "", "", "", err
	}
	// fmt.Printf("RC4 File Decryption Time: %.6f ms\n", float64(time.Since(startRC4).Nanoseconds())/1e6)
	// fmt.Printf("RC4 Decrypted File Size: %d bytes\n\n", len(decryptedRC4))

	decryptedRC4Path := filepath.Join(rc4DecryptedFolder, removeExtension(filepath.Base(rc4FilePath), ".rc4"))
	if err := os.WriteFile(decryptedRC4Path, decryptedRC4, 0644); err != nil {
		return "", "", "", err
	}

	// Decrypt and save DES file
	desDecryptedFolder := filepath.Join(decryptedPath, "des")
	if err := os.MkdirAll(desDecryptedFolder, 0755); err != nil {
		return "", "", "", err
	}
	desEncryptedData, err := os.ReadFile(desFilePath)
	if err != nil {
		return "", "", "", err
	}
	// startDES := time.Now()
	decryptedDES, err := DecryptFileBytesDES(desEncryptedData, secretKey8Byte)
	if err != nil {
		return "", "", "", err
	}
	// fmt.Printf("DES File Decryption Time: %.6f ms\n\n", float64(time.Since(startDES).Nanoseconds())/1e6)
	// fmt.Printf("DES Decrypted File Size: %d bytes\n\n", len(decryptedDES))

	decryptedDESPath := filepath.Join(desDecryptedFolder, removeExtension(filepath.Base(desFilePath), ".des"))
	if err := os.WriteFile(decryptedDESPath, decryptedDES, 0644); err != nil {
		return "", "", "", err
	}

	return decryptedAESPath, decryptedRC4Path, decryptedDESPath, nil
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
