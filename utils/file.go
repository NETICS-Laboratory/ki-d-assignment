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

func UploadFile(fileHeader *multipart.FileHeader, filePath string, secretKey []byte, secretKey8Byte []byte) error {
	// Buat direktori utama untuk user berdasarkan filePath (user.ID)
	if err := os.MkdirAll(filePath, 0755); err != nil {
		return err
	}

	// Buka file yang diunggah
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	// Baca file sebagai byte array
	fileData, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// Enkripsi file menggunakan AES, RC4, dan DES
	encryptedFileAES, err := EncryptFileBytesAES(fileData, secretKey)
	if err != nil {
		return err
	}

	encryptedFileRC4, err := EncryptFileBytesRC4(fileData, secretKey)
	if err != nil {
		return err
	}

	encryptedFileDES, err := EncryptFileBytesDES(fileData, secretKey8Byte)
	if err != nil {
		return err
	}

	// Simpan file yang dienkripsi di masing-masing direktori sesuai algoritma enkripsi

	// Simpan file terenkripsi AES di /aes
	aesPath := filepath.Join(filePath, "aes")
	if err := os.MkdirAll(aesPath, 0755); err != nil {
		return err
	}
	aesFile := filepath.Join(aesPath, fileHeader.Filename)
	if err := os.WriteFile(aesFile, encryptedFileAES, 0644); err != nil {
		return err
	}

	// Simpan file terenkripsi RC4 di /rc4
	rc4Path := filepath.Join(filePath, "rc4")
	if err := os.MkdirAll(rc4Path, 0755); err != nil {
		return err
	}
	rc4File := filepath.Join(rc4Path, fileHeader.Filename)
	if err := os.WriteFile(rc4File, encryptedFileRC4, 0644); err != nil {
		return err
	}

	// Simpan file terenkripsi DES di /des
	desPath := filepath.Join(filePath, "des")
	if err := os.MkdirAll(desPath, 0755); err != nil {
		return err
	}
	desFile := filepath.Join(desPath, fileHeader.Filename)
	if err := os.WriteFile(desFile, encryptedFileDES, 0644); err != nil {
		return err
	}

	return nil
}
