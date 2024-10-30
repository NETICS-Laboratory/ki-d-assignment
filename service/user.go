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
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type UserService interface {
	RegisterUser(ctx context.Context, userDTO dto.UserCreateDto) (entity.User, error)
	// GetAllUser(ctx context.Context) ([]entity.User, error)
	FindUserByUsername(ctx context.Context, username string) (entity.User, error)
	Verify(ctx context.Context, email string, password string) (bool, error)
	CheckUser(ctx context.Context, email string) (bool, error)
	// DeleteUser(ctx context.Context, userID uuid.UUID) error
	// UpdateUser(ctx context.Context, userDTO dto.UserUpdateDto) error
	MeUser(ctx context.Context, userID uuid.UUID) (entity.User, error)
	MeUserDecrypted(ctx context.Context, userID uuid.UUID) (dto.UserRequestDecryptedDto, error)
	DecryptUserIDCard(ctx context.Context, userID uuid.UUID) error
	RequestAccess(ctx context.Context, userID, allowedUserID uuid.UUID) (entity.AccessRequest, error)
	GetAccessRequests(ctx context.Context, userID uuid.UUID) ([]entity.AccessRequest, error)
	UpdateAccessRequestStatus(ctx context.Context, requestID uuid.UUID, status string) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(ur repository.UserRepository) UserService {
	return &userService{
		userRepository: ur,
	}
}

func (us *userService) RegisterUser(ctx context.Context, userDTO dto.UserCreateDto) (entity.User, error) {
	user := entity.User{}

	user.ID = uuid.New()

	// CREDENTIALS
	user.Username = userDTO.Username
	user.Username_AES = userDTO.Username
	user.Username_RC4 = userDTO.Username
	user.Username_DES = userDTO.Username

	user.Password_AES = userDTO.Password
	user.Password_RC4 = userDTO.Password
	user.Password_DES = userDTO.Password

	user.Name_AES = userDTO.Name
	user.Name_DES = userDTO.Name
	user.Name_RC4 = userDTO.Name

	user.Email_AES = userDTO.Email
	user.Email_DES = userDTO.Email
	user.Email_RC4 = userDTO.Email

	user.NoTelp_AES = userDTO.NoTelp
	user.NoTelp_DES = userDTO.NoTelp
	user.NoTelp_RC4 = userDTO.NoTelp

	user.Address_AES = userDTO.Address
	user.Address_DES = userDTO.Address
	user.Address_RC4 = userDTO.Address

	// KEYS
	secretKey, err := utils.GenerateSecretKey()
	if err != nil {
		return entity.User{}, err
	}
	user.SecretKey = secretKey

	secretKey8Byte, err := utils.GenerateSecretKey8Byte()
	if err != nil {
		return entity.User{}, err
	}
	user.SecretKey8Byte = secretKey8Byte

	// Generate asymmetric keys for the user
	if err := utils.GenerateAsymmetricKeys(user.ID); err != nil {
		return entity.User{}, err
	}

	// ID CARD
	// Validate file type for id card
	fileHeader := userDTO.ID_Card
	file, err := fileHeader.Open()
	if err != nil {
		return entity.User{}, err
	}
	defer file.Close()

	fileType := fileHeader.Header.Get("Content-Type")
	if fileType != "image/jpeg" && fileType != "image/png" {
		return entity.User{}, errors.New("unsupported file type: only JPEG and PNG are allowed")
	}

	// Validate file name
	if !utils.IsValidFileName(fileHeader.Filename) {
		return entity.User{}, errors.New("invalid file name")
	}

	user.ID_Card_ID = uuid.New()

	// Upload and encrypt the ID Card
	filePath := fmt.Sprintf("uploads/%s/encrypted", user.Username)
	aesFile, rc4File, desFile, err := utils.UploadFile(fileHeader, filePath, user.SecretKey, user.SecretKey8Byte)
	if err != nil {
		return entity.User{}, err
	}

	// Store the paths of the encrypted files in the entity
	user.ID_Card_AES = aesFile
	user.ID_Card_RC4 = rc4File
	user.ID_Card_DES = desFile

	// fmt.Printf("%v\n%v\n%v", user.ID_Card_AES, user.ID_Card_RC4, user.ID_Card_DES)

	return us.userRepository.RegisterUser(ctx, user)
}

// func (us *userService) GetAllUser(ctx context.Context) ([]entity.User, error) {
// 	return us.userRepository.GetAllUser(ctx)
// }

func (us *userService) FindUserByUsername(ctx context.Context, username string) (entity.User, error) {
	return us.userRepository.FindUserByUsername(ctx, username)
}

func (us *userService) Verify(ctx context.Context, username string, password string) (bool, error) {
	res, err := us.userRepository.FindUserByUsername(ctx, username)
	if err != nil {
		return false, err
	}

	decrypted, err := helpers.DecryptData(res.Password_AES, res.Password_RC4, res.Password_DES, res.SecretKey, res.SecretKey8Byte)
	if err != nil {
		return false, err
	}

	if decrypted == password {
		return true, nil
	}

	return false, nil
}

func (us *userService) CheckUser(ctx context.Context, username string) (bool, error) {
	result, err := us.userRepository.FindUserByUsername(ctx, username)
	if err != nil {
		return false, err
	}

	if result.Username == "" {
		return false, nil
	}
	return true, nil
}

func (us *userService) MeUser(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	return us.userRepository.FindUserByID(ctx, userID)
}

func (us *userService) MeUserDecrypted(ctx context.Context, userID uuid.UUID) (dto.UserRequestDecryptedDto, error) {
	res, err := us.userRepository.FindUserByID(ctx, userID)
	if err != nil {
		return dto.UserRequestDecryptedDto{}, err
	}

	dataToDecrypt := []struct {
		aesCipher string
		rc4Cipher string
		desCipher string
		decrypted *string
	}{
		{res.Username_AES, res.Username_RC4, res.Username_DES, new(string)},
		{res.Name_AES, res.Name_RC4, res.Name_DES, new(string)},
		{res.Email_AES, res.Email_RC4, res.Email_DES, new(string)},
		{res.NoTelp_AES, res.NoTelp_RC4, res.NoTelp_DES, new(string)},
		{res.Address_AES, res.Address_RC4, res.Address_DES, new(string)},
	}

	for _, data := range dataToDecrypt {
		decrypted, err := helpers.DecryptData(data.aesCipher, data.rc4Cipher, data.desCipher, res.SecretKey, res.SecretKey8Byte)
		if err != nil {
			return dto.UserRequestDecryptedDto{}, err
		}
		*data.decrypted = decrypted
	}

	id_card_aes, id_card_rc4, id_card_des, err := helpers.DecryptDataReturnIndiviual(res.ID_Card_AES, res.ID_Card_RC4, res.ID_Card_DES, res.SecretKey, res.SecretKey8Byte)
	if err != nil {
		return dto.UserRequestDecryptedDto{}, err
	}

	// baseEncryptedPath := filepath.Join("uploads", res.Username, "encrypted")

	normalizeFilePath := func(encryptedFilePath, encryptionType string) string {
		// Remove the encryption folder (aes/rc4/des)
		filePath := strings.ReplaceAll(encryptedFilePath, filepath.Join("encrypted", encryptionType), "encrypted")

		// Remove the file extension (aes, rc4, des)
		return strings.TrimSuffix(filePath, "."+encryptionType)
	}

	id_card_aes = normalizeFilePath(id_card_aes, "aes")
	id_card_rc4 = normalizeFilePath(id_card_rc4, "rc4")
	id_card_des = normalizeFilePath(id_card_des, "des")

	if id_card_aes != id_card_rc4 && id_card_rc4 != id_card_des {
		return dto.UserRequestDecryptedDto{}, errors.New("ID Card Did Not Same")
	}

	result := dto.UserRequestDecryptedDto{
		ID:       res.ID,
		Username: *dataToDecrypt[0].decrypted,
		Name:     *dataToDecrypt[1].decrypted,
		Email:    *dataToDecrypt[2].decrypted,
		NoTelp:   *dataToDecrypt[3].decrypted,
		Address:  *dataToDecrypt[4].decrypted,
		ID_Card:  id_card_aes,
	}

	return result, nil
}

func (us *userService) DecryptUserIDCard(ctx context.Context, userID uuid.UUID) error {
	// Find user by ID
	user, err := us.userRepository.FindUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Decrypt file paths for AES, RC4, and DES
	decryptedAESPath, decryptedRC4Path, decryptedDESPath, err := helpers.DecryptDataReturnIndiviual(user.ID_Card_AES, user.ID_Card_RC4, user.ID_Card_DES, user.SecretKey, user.SecretKey8Byte)
	if err != nil {
		return fmt.Errorf("gagal melakukan dekripi path id card: %v", err)
	}
	// fmt.Printf("%v %v %v", decryptedAESPath, decryptedRC4Path, decryptedDESPath)
	// Now use the decrypted file paths to call the DecryptAndSaveFiles function
	filePath := fmt.Sprintf("uploads/%s", user.Username)

	err = utils.DecryptAndSaveFiles(filePath, decryptedAESPath, decryptedRC4Path, decryptedDESPath, user.SecretKey, user.SecretKey8Byte)
	if err != nil {
		return fmt.Errorf("gagal melakukan dekripi id card: %v", err)
	}

	return nil
}

func (us *userService) RequestAccess(ctx context.Context, userID, allowedUserID uuid.UUID) (entity.AccessRequest, error) {
	request := entity.AccessRequest{
		UserID:        userID,
		AllowedUserID: allowedUserID,
		Status:        "pending",
	}
	return us.userRepository.CreateAccessRequest(ctx, request)
}

func (us *userService) GetAccessRequests(ctx context.Context, userID uuid.UUID) ([]entity.AccessRequest, error) {
	return us.userRepository.GetAccessRequestsByUserID(ctx, userID)
}

func (us *userService) UpdateAccessRequestStatus(ctx context.Context, requestID uuid.UUID, status string) error {
	return us.userRepository.UpdateAccessRequestStatus(ctx, requestID, status)
}
