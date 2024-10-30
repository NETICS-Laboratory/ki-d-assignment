package service

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"ki-d-assignment/dto"
	"ki-d-assignment/entity"
	"ki-d-assignment/helpers"
	"ki-d-assignment/repository"
	"ki-d-assignment/utils"
	"os"
	"path/filepath"
	"strings"
	"text/template"

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

	RequestAccess(ctx context.Context, userID, requestedUserID uuid.UUID) (entity.AccessRequest, error)
	GetSentAccessRequests(ctx context.Context, userID uuid.UUID) ([]entity.AccessRequest, error)
	GetReceivedAccessRequests(ctx context.Context, userID uuid.UUID) ([]entity.AccessRequest, error)
	UpdateAccessRequestStatus(ctx context.Context, userID, requestID uuid.UUID, status string) error
}

type userService struct {
	userRepository          repository.UserRepository
	accessRequestRepository repository.AccessRequestRepository
}

func NewUserService(ur repository.UserRepository, arr repository.AccessRequestRepository) UserService {
	return &userService{
		userRepository:          ur,
		accessRequestRepository: arr,
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
	if err := utils.GenerateAsymmetricKeys(user.Username); err != nil {
		return entity.User{}, fmt.Errorf("gagal membuat RSA key pair: %v", err)
	}

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

// func (us *userService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
// 	return us.userRepository.DeleteUser(ctx, userID)
// }

// func (us *userService) UpdateUser(ctx context.Context, userDTO dto.UserUpdateDto) error {
// 	user, err := us.userRepository.FindUserByID(ctx, userDTO.ID)
// 	if err != nil {
// 		return err
// 	}

// 	if userDTO.Name != "" {
// 		user.Name = userDTO.Name
// 	}
// 	if userDTO.Email != "" {
// 		user.Email = userDTO.Email
// 	}
// 	if userDTO.NoTelp != "" {
// 		user.NoTelp = userDTO.NoTelp
// 	}
// 	if userDTO.Password != "" {
// 		// Hash new password before saving
// 		hashedPassword, err := utils.HashPassword(userDTO.Password)
// 		if err != nil {
// 			return err
// 		}
// 		user.Password = hashedPassword
// 	}

// 	return us.userRepository.UpdateUser(ctx, user)
// }

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

func (us *userService) RequestAccess(ctx context.Context, userID, requestedUserID uuid.UUID) (entity.AccessRequest, error) {
	exists, err := us.accessRequestRepository.CheckExistingAccessRequest(ctx, userID, requestedUserID)
	if err != nil {
		return entity.AccessRequest{}, err
	}
	if exists {
		return entity.AccessRequest{}, errors.New("Akses Request Ke User Tersebut Sudah Ada")
	}

	request := entity.AccessRequest{
		UserID:          userID,
		RequestedUserID: requestedUserID,
		Status:          "pending",
	}
	return us.accessRequestRepository.CreateAccessRequest(ctx, request)
}

func (us *userService) GetSentAccessRequests(ctx context.Context, userID uuid.UUID) ([]entity.AccessRequest, error) {
	return us.accessRequestRepository.GetAccessRequestsByUserID(ctx, userID)
}

func (us *userService) GetReceivedAccessRequests(ctx context.Context, userID uuid.UUID) ([]entity.AccessRequest, error) {
	return us.accessRequestRepository.GetAccessRequestsByRequestedUserID(ctx, userID)
}

func (us *userService) UpdateAccessRequestStatus(ctx context.Context, userID, requestID uuid.UUID, status string) error {
	res, err := us.accessRequestRepository.GetAccessRequestsByID(ctx, requestID)
	if err != nil {
		return err
	}

	if res.RequestedUserID != userID {
		return fmt.Errorf("Anda tidak memiliki izin untuk mengubah status permintaan ini")
	}

	if err := us.accessRequestRepository.UpdateAccessRequestStatus(ctx, requestID, status); err != nil {
		return err
	}

	if status == "approved" {
		// Ambil requesting user ID dari requestnya
		requestingUser, err := us.userRepository.FindUserByID(ctx, res.UserID)
		if err != nil {
			return err
		}

		// Ambil public key si requesting user
		publicKeyPath := filepath.Join("uploads", requestingUser.Username, "secret", "public_key.pem")
		publicKeyFile, err := os.ReadFile(publicKeyPath)
		if err != nil {
			return err
		}

		block, _ := pem.Decode(publicKeyFile)
		if block == nil || block.Type != "RSA PUBLIC KEY" {
			return errors.New("failed to decode public key")
		}

		publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return err
		}

		// Encrypt secret key requested user with User A's public key

		requestedUser, err := us.userRepository.FindUserByID(ctx, userID)
		if err != nil {
			return err
		}

		encryptedSecretKey, err := helpers.EncryptWithPublicKey(requestedUser.SecretKey, publicKey)
		if err != nil {
			return err
		}

		encryptedSecretKey8Byte, err := helpers.EncryptWithPublicKey(requestedUser.SecretKey8Byte, publicKey)
		if err != nil {
			return err
		}

		decryptedRequestingUserEmail, err := helpers.DecryptData(requestingUser.Email_AES, requestingUser.Email_RC4, requestingUser.Email_DES, requestingUser.SecretKey, requestingUser.SecretKey8Byte)
		if err != nil {
			return fmt.Errorf("failed to decrypt requesting user's email: %v", err)
		}

		emailBody, err := buildKeyAccessEmail(
			requestingUser.Username,
			requestedUser.Username,
			fmt.Sprintf("%x", encryptedSecretKey),
			fmt.Sprintf("%x", encryptedSecretKey8Byte),
		)
		if err != nil {
			return fmt.Errorf("failed to build email body: %v", err)
		}

		if err := helpers.SendMail(decryptedRequestingUserEmail, "Access Approved", emailBody); err != nil {
			return err
		}
	}

	return nil
}

func buildKeyAccessEmail(username, requestedUsername, encryptedSecretKey, encryptedSecretKey8Byte string) (string, error) {
	// Load the HTML template
	htmlFilePath := "helpers/email-template/access-approved.html"
	htmlTemplate, err := ioutil.ReadFile(htmlFilePath)
	if err != nil {
		return "", err
	}

	// Parse the HTML template with placeholders
	tmpl, err := template.New("email").Parse(string(htmlTemplate))
	if err != nil {
		return "", err
	}

	// Data to inject into the template
	data := struct {
		Username                string
		RequestedUsername       string
		EncryptedSecretKey      string
		EncryptedSecretKey8Byte string
	}{
		Username:                username,
		RequestedUsername:       requestedUsername,
		EncryptedSecretKey:      encryptedSecretKey,
		EncryptedSecretKey8Byte: encryptedSecretKey8Byte,
	}

	// Execute the template with the data
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return "", err
	}

	return body.String(), nil
}
