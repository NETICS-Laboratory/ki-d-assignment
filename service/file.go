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

	"github.com/google/uuid"
)

type FileService interface {
	UploadFile(ctx context.Context, fileDTO dto.FileCreateDto, userID uuid.UUID) (entity.Files, error)
	GetUserFiles(ctx context.Context, userID uuid.UUID) ([]entity.Files, error)
	GetUserFileDecryptedByID(ctx context.Context, fileID uuid.UUID, userID uuid.UUID) (dto.FileDecryptedResponse, error)
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
