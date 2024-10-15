package service

import (
	"context"
	"errors"
	"fmt"
	"ki-d-assignment/dto"
	"ki-d-assignment/entity"
	"ki-d-assignment/repository"
	"ki-d-assignment/utils"

	"github.com/google/uuid"
)

type FileService interface {
	UploadFile(ctx context.Context, fileDTO dto.FileCreateDto, userID uuid.UUID) (entity.Files, error)
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
