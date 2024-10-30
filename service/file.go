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

	"github.com/google/uuid"
)

type FileService interface {
	UploadFile(ctx context.Context, fileDTO dto.FileCreateDto, userID uuid.UUID) (entity.Files, error)
	GetUserFiles(ctx context.Context, userID uuid.UUID) ([]entity.Files, error)
	GetUserFileDecryptedByID(ctx context.Context, fileID uuid.UUID, userID uuid.UUID) (dto.FileDecryptedResponse, error)
	GetRequestedUserFile(ctx context.Context, userID uuid.UUID, requestedUser entity.User, secretKeys string, secretKeys8Byte string) ([]dto.FileDecryptedResponse, error)
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
	if (err != nil) {
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

func (fs *fileService) GetRequestedUserFile(ctx context.Context, userID uuid.UUID, requestedUser entity.User, secretKeys string, secretKeys8Byte string) ([]dto.FileDecryptedResponse, error) {
	// Retrieve files for the requested user
	files, err := fs.fileRepository.GetFilesByUserID(ctx, requestedUser.ID)
	if err != nil {
		return nil, fmt.Errorf("file tidak ditemukan atau tidak memiliki akses: %v", err)
	}

	// Fetch the user who requested access
	user, err := fs.userRepository.FindUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal menemukan pengguna: %v", err)
	}

	// Decode the hex-encoded secret keys
	decodedSecretKey, err := hex.DecodeString(secretKeys)
	if err != nil {
		return nil, fmt.Errorf("failed to decode secret key: %v", err)
	}

	decodedSecretKey8Byte, err := hex.DecodeString(secretKeys8Byte)
	if err != nil {
		return nil, fmt.Errorf("failed to decode 8-byte secret key: %v", err)
	}

	// Validate key lengths
	if len(decodedSecretKey) != 32 { // Assuming AES-256 and RC4-256
		return nil, fmt.Errorf("invalid AES/RC4 key length: expected 32 bytes, got %d bytes", len(decodedSecretKey))
	}

	if len(decodedSecretKey8Byte) != 8 { // DES requires 8-byte key
		return nil, fmt.Errorf("invalid DES key length: expected 8 bytes, got %d bytes", len(decodedSecretKey8Byte))
	}

	// Define the file path for saving decrypted files
	filePath := fmt.Sprintf("uploads/%s", user.Username)

	// Array to hold decrypted responses
	var decryptedResponses []dto.FileDecryptedResponse

	// Decrypt each file and build the response
	for _, file := range files {
		// Decrypt AES file content
		decryptedAES, err := utils.DecryptFileBytesAES([]byte(file.Files_AES), decodedSecretKey)
		if err != nil {
			return nil, fmt.Errorf("gagal melakukan dekripsi AES: %v", err)
		}

		// Decrypt RC4 file content
		decryptedRC4, err := utils.DecryptFileBytesRC4([]byte(file.Files_RC4), decodedSecretKey)
		if err != nil {
			return nil, fmt.Errorf("gagal melakukan dekripsi RC4: %v", err)
		}

		// Decrypt DES file content
		decryptedDES, err := utils.DecryptFileBytesDES([]byte(file.Files_DES), decodedSecretKey8Byte)
		if err != nil {
			return nil, fmt.Errorf("gagal melakukan dekripsi DES: %v", err)
		}

		// Save decrypted files to disk
		err = utils.DecryptAndSaveFiles(filePath, string(decryptedAES), string(decryptedRC4), string(decryptedDES), decodedSecretKey, decodedSecretKey8Byte)
		if err != nil {
			return nil, fmt.Errorf("gagal menyimpan file yang telah didekripsi: %v", err)
		}

		// Append decrypted paths to the response array
		decryptedResponses = append(decryptedResponses, dto.FileDecryptedResponse{
			ID:            file.ID,
			Decrypted_AES: string(decryptedAES),
			Decrypted_RC4: string(decryptedRC4),
			Decrypted_DES: string(decryptedDES),
		})
	}

	return decryptedResponses, nil
}

