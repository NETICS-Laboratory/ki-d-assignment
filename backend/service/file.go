package service

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"ki-d-assignment/dto"
	"ki-d-assignment/entity"
	"ki-d-assignment/helpers"
	"ki-d-assignment/repository"
	"ki-d-assignment/utils"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type FileService interface {
	UploadFile(ctx context.Context, fileDTO dto.FileCreateDto, userID uuid.UUID) (entity.Files, error)
	GetUserFiles(ctx context.Context, userID uuid.UUID) ([]entity.Files, error)
	GetUserFileDecryptedByID(ctx context.Context, fileID uuid.UUID, userID uuid.UUID) (dto.FileDecryptedResponse, error)
	GetRequestedUserData(ctx context.Context, requestedUser entity.User, secretKeys string, secretKeys8Byte string) ([]dto.FileDecryptedResponse, error)
	SignPDF(ctx context.Context, fileID uuid.UUID, userID uuid.UUID) (entity.Files, error)
}

type fileService struct {
	fileRepository repository.FileRepository
	userRepository repository.UserRepository
}

func NewFileService(fr repository.FileRepository, ur repository.UserRepository) FileService {
	return &fileService{
		fileRepository: fr,
		userRepository: ur,
	}
}

func (fs *fileService) UploadFile(ctx context.Context, fileDTO dto.FileCreateDto, userID uuid.UUID) (entity.Files, error) {

	user, err := fs.userRepository.FindUserByID(ctx, userID)
	if err != nil {
		return entity.Files{}, fmt.Errorf("failed to retrieve user data: %v", err)
	}

	file := entity.Files{
		ID:     uuid.New(),
		UserID: user.ID,
	}

	// Validate file type for files (can be PDF/DOC/XLS files, and video files)
	fileType := fileDTO.File.Header.Get("Content-Type")
	if fileType != "application/pdf" &&
		fileType != "application/msword" &&
		fileType != "application/vnd.openxmlformats-officedocument.wordprocessingml.document" && // Untuk .docx
		fileType != "application/vnd.ms-excel" &&
		fileType != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" && // Untuk .xlsx
		fileType != "video/mp4" && fileType != "image/jpeg" && fileType != "image/png" {
		return entity.Files{}, errors.New("unsupported file type")
	}

	// Validate file name
	if !utils.IsValidFileName(fileDTO.File.Filename) {
		return entity.Files{}, errors.New("invalid file name")
	}
	// Upload and encrypt the file path
	encryptedFilePath := fmt.Sprintf("uploads/%s/encrypted", user.Username)

	secretKey := user.SecretKey
	secretKey8Byte := user.SecretKey8Byte

	aesFile, rc4File, desFile, err := utils.UploadFile(fileDTO.File, encryptedFilePath, secretKey, secretKey8Byte)
	if err != nil {
		return entity.Files{}, fmt.Errorf("failed to upload and encrypt file: %v", err)
	}
	// Store the paths of the encrypted files in the entity

	file.Files_AES = aesFile
	file.Files_RC4 = rc4File
	file.Files_DES = desFile

	uploadedFile, err := fs.fileRepository.UploadFile(ctx, file)
	if err != nil {
		return entity.Files{}, err
	}
	return uploadedFile, nil
}

func (fs *fileService) GetUserFiles(ctx context.Context, userID uuid.UUID) ([]entity.Files, error) {
	files, err := fs.fileRepository.GetFilesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (fs *fileService) GetUserFileDecryptedByID(ctx context.Context, fileID uuid.UUID, userID uuid.UUID) (dto.FileDecryptedResponse, error) {
	// 1. Find the user by ID to get the secret keys
	file, err := fs.fileRepository.GetFileByIDAndUserID(ctx, fileID, userID)
	if err != nil {
		return dto.FileDecryptedResponse{}, fmt.Errorf("file tidak ditemukan atau tidak memiliki akses: %v", err)
	}

	// 2. Find the specific file by ID
	user, err := fs.userRepository.FindUserByID(ctx, userID)
	if err != nil {
		return dto.FileDecryptedResponse{}, fmt.Errorf("gagal menemukan pengguna: %v", err)
	}

	// 3. Get the file paths (AES, RC4, DES) from the file entity
	decryptedAESPath, decryptedRC4Path, decryptedDESPath, err := helpers.DecryptDataReturnIndiviual(
		file.Files_AES, file.Files_RC4, file.Files_DES, user.SecretKey, user.SecretKey8Byte)
	if err != nil {
		return dto.FileDecryptedResponse{}, fmt.Errorf("gagal melakukan dekripsi jalur file: %v", err)
	}

	// 4. Call the utility function to decrypt the files and save them in the "decrypted" folder
	filePath := fmt.Sprintf("uploads/%s", user.Username)
	err = utils.DecryptAndSaveFiles(filePath, decryptedAESPath, decryptedRC4Path, decryptedDESPath, user.SecretKey, user.SecretKey8Byte)
	if err != nil {
		return dto.FileDecryptedResponse{}, fmt.Errorf("gagal melakukan dekripsi file dan menyimpannya: %v", err)
	}

	decryptedResponse := dto.FileDecryptedResponse{
		ID:            file.ID,
		Decrypted_AES: decryptedAESPath,
		Decrypted_RC4: decryptedRC4Path,
		Decrypted_DES: decryptedDESPath,
	}

	return decryptedResponse, nil
}

// TODO: Add more validation
func (fs *fileService) GetRequestedUserData(ctx context.Context, requestedUser entity.User, secretKeys string, secretKeys8Byte string) ([]dto.FileDecryptedResponse, error) {
	files, err := fs.fileRepository.GetFilesByUserID(ctx, requestedUser.ID)
	if err != nil {
		return nil, fmt.Errorf("file tidak ditemukan atau tidak memiliki akses: %v", err)
	}

	decodedSecretKey, err := hex.DecodeString(secretKeys)
	if err != nil {
		return nil, fmt.Errorf("failed to decode secret key: %v", err)
	}

	decodedSecretKey8Byte, err := hex.DecodeString(secretKeys8Byte)
	if err != nil {
		return nil, fmt.Errorf("failed to decode 8-byte secret key: %v", err)
	}

	if len(decodedSecretKey) != 32 { // AES-256 and RC4-256
		return nil, fmt.Errorf("invalid AES/RC4 key length: expected 32 bytes, got %d bytes", len(decodedSecretKey))
	}

	if len(decodedSecretKey8Byte) != 8 { // DES 8-byte key
		return nil, fmt.Errorf("invalid DES key length: expected 8 bytes, got %d bytes", len(decodedSecretKey8Byte))
	}

	filePath := fmt.Sprintf("uploads/%s", requestedUser.Username)

	var decryptedResponses []dto.FileDecryptedResponse

	for _, file := range files {

		decryptedAES, decryptedRC4, decryptedDES, err := helpers.DecryptDataReturnIndiviual(file.Files_AES, file.Files_RC4, file.Files_DES, decodedSecretKey, decodedSecretKey8Byte)
		if err != nil {
			return nil, fmt.Errorf("gagal melakukan dekripsi file path: %v", err)
		}

		err = utils.DecryptAndSaveFiles(filePath, decryptedAES, decryptedRC4, decryptedDES, decodedSecretKey, decodedSecretKey8Byte)
		if err != nil {
			return nil, fmt.Errorf("gagal menyimpan file yang telah didekripsi: %v", err)
		}

		decryptedResponses = append(decryptedResponses, dto.FileDecryptedResponse{
			ID:            file.ID,
			Decrypted_AES: decryptedAES,
			Decrypted_RC4: decryptedRC4,
			Decrypted_DES: decryptedDES,
		})
	}

	return decryptedResponses, nil
}

func (fs *fileService) SignPDF(ctx context.Context, fileID uuid.UUID, userID uuid.UUID) (entity.Files, error) {
	file, err := fs.fileRepository.GetFileByIDAndUserID(ctx, fileID, userID)
	if err != nil {
		return entity.Files{}, fmt.Errorf("file not found or access denied: %v", err)
	}

	user, err := fs.userRepository.FindUserByID(ctx, userID)
	if err != nil {
		return entity.Files{}, fmt.Errorf("user not found: %v", err)
	}

	decryptedAESPath, decryptedRC4Path, decryptedDESPath, err := helpers.DecryptDataReturnIndiviual(
		file.Files_AES, file.Files_RC4, file.Files_DES, user.SecretKey, user.SecretKey8Byte)
	if err != nil {
		return entity.Files{}, fmt.Errorf("gagal melakukan dekripsi jalur file: %v", err)
	}

	filePath := fmt.Sprintf("uploads/%s", user.Username)
	AESPath, _, _, err := utils.DecryptAndSaveFilesReturnPath(filePath, decryptedAESPath, decryptedRC4Path, decryptedDESPath, user.SecretKey, user.SecretKey8Byte)
	if err != nil {
		return entity.Files{}, fmt.Errorf("gagal melakukan dekripsi file dan menyimpannya: %v", err)
	}

	decryptedFilePath := fmt.Sprintf("uploads/%s/decrypted/aes/%s", user.Username, filepath.Base(AESPath))
	fmt.Printf("Looking for decrypted file at: %s\n", decryptedFilePath)

	privateKeyPath := fmt.Sprintf("uploads/%s/secret/private_key.pem", user.Username)
	fmt.Printf("Looking for private key at: %s\n", privateKeyPath)

	signedFilePath := fmt.Sprintf("uploads/%s/signed/%s.signed", user.Username, filepath.Base(AESPath))
	fmt.Printf("Will save signed file at: %s\n", signedFilePath)

	// Check if decrypted file exists
	if _, err := os.Stat(decryptedFilePath); os.IsNotExist(err) {
		return entity.Files{}, fmt.Errorf("decrypted file not found: %s", decryptedFilePath)
	}

	// Check if private key exists
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		return entity.Files{}, fmt.Errorf("private key not found: %s", privateKeyPath)
	}

	// Sign the file
	err = utils.SignPDFWithOpenSSL(decryptedFilePath, privateKeyPath, signedFilePath)
	if err != nil {
		return entity.Files{}, fmt.Errorf("failed to sign PDF: %v", err)
	}

	// Update the file record with the signed file path
	file.Signature = signedFilePath
	if err := fs.fileRepository.UpdateFile(file); err != nil {
		return entity.Files{}, fmt.Errorf("failed to update file record: %v", err)
	}

	return file, nil
}