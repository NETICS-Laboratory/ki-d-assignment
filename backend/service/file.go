package service

import (
	"context"
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
	GetFileByID(fileID uuid.UUID, userID uuid.UUID) (entity.Files, error)
	UpdateFile(ctx context.Context, file entity.Files) error
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

	fileType := fileDTO.File.Header.Get("Content-Type")
	validTypes := map[string]bool{
		"application/pdf": true, "application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/vnd.ms-excel": true, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
		"video/mp4": true, "image/jpeg": true, "image/png": true,
	}
	if !validTypes[fileType] {
		return entity.Files{}, errors.New("unsupported file type")
	}

	if !utils.IsValidFileName(fileDTO.File.Filename) {
		return entity.Files{}, errors.New("invalid file name")
	}

	encryptedFilePath := fmt.Sprintf("uploads/%s/encrypted", user.Username)
	aesFile, rc4File, desFile, err := utils.UploadFile(fileDTO.File, encryptedFilePath, user.SecretKey, user.SecretKey8Byte)
	if err != nil {
		return entity.Files{}, fmt.Errorf("failed to upload and encrypt file: %v", err)
	}

	file.Files_AES = aesFile
	file.Files_RC4 = rc4File
	file.Files_DES = desFile

	return fs.fileRepository.UploadFile(ctx, file)
}

func (fs *fileService) GetUserFiles(ctx context.Context, userID uuid.UUID) ([]entity.Files, error) {
	return fs.fileRepository.GetFilesByUserID(ctx, userID)
}

func (fs *fileService) GetUserFileDecryptedByID(ctx context.Context, fileID uuid.UUID, userID uuid.UUID) (dto.FileDecryptedResponse, error) {
	file, err := fs.fileRepository.GetFileByIDAndUserID(ctx, fileID, userID)
	if err != nil {
		return dto.FileDecryptedResponse{}, fmt.Errorf("file not found or access denied: %v", err)
	}

	user, err := fs.userRepository.FindUserByID(ctx, userID)
	if err != nil {
		return dto.FileDecryptedResponse{}, fmt.Errorf("failed to retrieve user: %v", err)
	}

	decryptedAESPath, decryptedRC4Path, decryptedDESPath, err := helpers.DecryptDataReturnIndiviual(
		file.Files_AES, file.Files_RC4, file.Files_DES, user.SecretKey, user.SecretKey8Byte)
	if err != nil {
		return dto.FileDecryptedResponse{}, fmt.Errorf("failed to decrypt file paths: %v", err)
	}

	decryptedFolderPath := fmt.Sprintf("uploads/%s/decrypted", user.Username)
	if err = utils.DecryptAndSaveFiles(decryptedFolderPath, decryptedAESPath, decryptedRC4Path, decryptedDESPath, user.SecretKey, user.SecretKey8Byte); err != nil {
		return dto.FileDecryptedResponse{}, fmt.Errorf("failed to save decrypted files: %v", err)
	}

	return dto.FileDecryptedResponse{
		ID:            file.ID,
		Decrypted_AES: decryptedAESPath,
		Decrypted_RC4: decryptedRC4Path,
		Decrypted_DES: decryptedDESPath,
	}, nil
}

func (fs *fileService) GetRequestedUserData(ctx context.Context, requestedUser entity.User, secretKeys string, secretKeys8Byte string) ([]dto.FileDecryptedResponse, error) {
	files, err := fs.fileRepository.GetFilesByUserID(ctx, requestedUser.ID)
	if err != nil {
		return nil, fmt.Errorf("file not found or access denied: %v", err)
	}

	var responses []dto.FileDecryptedResponse
	for _, file := range files {
		decryptedAES, decryptedRC4, decryptedDES, err := helpers.DecryptDataReturnIndiviual(file.Files_AES, file.Files_RC4, file.Files_DES, []byte(secretKeys), []byte(secretKeys8Byte))
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt file paths: %v", err)
		}
		responses = append(responses, dto.FileDecryptedResponse{
			ID:            file.ID,
			Decrypted_AES: decryptedAES,
			Decrypted_RC4: decryptedRC4,
			Decrypted_DES: decryptedDES,
		})
	}
	return responses, nil
}

func (fs *fileService) GetFileByID(fileID uuid.UUID, userID uuid.UUID) (entity.Files, error) {
	return fs.fileRepository.GetFileByIDAndUserID(context.Background(), fileID, userID)
}

func (fs *fileService) UpdateFile(ctx context.Context, file entity.Files) error {
	return fs.fileRepository.UpdateFile(file)
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

	decryptedFilePath := fmt.Sprintf("uploads/%s/decrypted/decrypted/aes/%s", user.Username, filepath.Base(file.Files_AES))
	fmt.Printf("Looking for decrypted file at: %s\n", decryptedFilePath)

	privateKeyPath := fmt.Sprintf("uploads/%s/secret/private_key.pem", user.Username)
	fmt.Printf("Looking for private key at: %s\n", privateKeyPath)

	signedFilePath := fmt.Sprintf("uploads/%s/signed/%s.signed", user.Username, filepath.Base(file.Files_AES))
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
