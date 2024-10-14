package service

import (
	"context"
	"errors"
	"fmt"
	"ki-d-assignment/dto"
	"ki-d-assignment/entity"
	"ki-d-assignment/helpers"
	"ki-d-assignment/repository"

	"github.com/google/uuid"
)

type FileService interface {
	UploadFile(ctx context.Context, fileDTO dto.FileCreateDto) (entity.Files, error)
	GetAllFiles(ctx context.Context) ([]entity.Files, error)
	DecryptFile(filename string, encryptionMethod string) (string, error)
	GetFile(ctx context.Context, fileID string) (entity.Files, error)
	GetFileByUserID(ctx context.Context, userID string) ([]entity.Files, error)
}

type fileService struct {
	FileRepository repository.FileRepository
}

func NewFileService(fileRepo repository.FileRepository) FileService {
	return &fileService{
		FileRepository: fileRepo,
	}
}

func (f *fileService) UploadFile(ctx context.Context, fileDTO dto.FileCreateDto) (entity.Files, error) {
	var file entity.Files

	file.ID = uuid.New()
	file.Name = fileDTO.Name
	file.Files_AES = fileDTO.Files.Filename
	file.Files_RC4 = fileDTO.Files.Filename
	file.Files_DEC = fileDTO.Files.Filename
	file.UserID, _ = uuid.Parse(fileDTO.UserID)

	// Check file type
	if fileDTO.Files.Header.Get("Content-Type") != "application/pdf" && fileDTO.Files.Header.Get("Content-Type") != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" && fileDTO.Files.Header.Get("Content-Type") != "application/vnd.openxmlformats-officedocument.wordprocessingml.document" && fileDTO.Files.Header.Get("Content-Type") != "image/jpeg" && fileDTO.Files.Header.Get("Content-Type") != "image/png" && fileDTO.Files.Header.Get("Content-Type") != "video/mp4" {
		return entity.Files{}, errors.New("file type is not supported")
	}

	// Check file size
	if fileDTO.Files.Size > 1000000 {
		return entity.Files{}, errors.New("file size is too large")
	}

	// Check file name
	if fileDTO.Files.Filename == "" {
		return entity.Files{}, errors.New("file name is not valid")
	}

	// Save the files to the uploads folder
	fileName := fmt.Sprintf("%s/files/%s", file.UserID, file.ID)
	if err := utils.UploadFileUtility(fileDTO.Files, fileName); err != nil {
		return entity.Files{}, err
	}

	result, err := f.FileRepository.UploadFile(ctx, file)
	if err != nil {
		return entity.Files{}, err
	}

	return result, nil
}

func (f *fileService) GetAllFiles(ctx context.Context) ([]entity.Files, error) {
	result, err := f.FileRepository.GetAllFiles(ctx)
	if err != nil {
		return []entity.Files{}, err
	}

	return result, nil
}

func (f *fileService) DecryptFile(filename string, encryptionMethod string) (string, error) {

	if encryptionMethod == "AES" {
		return utils.DecryptAES(filename)
	} else if encryptionMethod == "RC4" {
		return utils.DecryptRC4(filename)
	} else if encryptionMethod == "DES" {
		return utils.DecryptDES(filename)
	} else {
		return "", errors.New("encryption method is not valid")
	}
}

// Get File from repository
func (f *fileService) GetFile(ctx context.Context, fileID string) (entity.Files, error) {
	result, err := f.FileRepository.GetFile(ctx, fileID)
	if err != nil {
		return entity.Files{}, err
	}

	return result, nil
}

// Get File by User id
func (f *fileService) GetFileByUserID(ctx context.Context, userID string) ([]entity.Files, error) {
	result, err := f.FileRepository.GetFileByUserID(ctx, userID)
	if err != nil {
		return []entity.Files{}, err
	}

	return result, nil
}