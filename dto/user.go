package dto

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type UserCreateDto struct {
	ID       uuid.UUID             `gorm:"type:char(36);primary_key;not_null" json:"id"`
	Name     string                `json:"name" form:"name" binding:"required"`
	Email    string                `json:"email" form:"email" binding:"required,email"`
	NoTelp   string                `json:"no_telp" form:"no_telp" binding:"required"`
	Address  string                `json:"address" form:"address" binding:"required"` // Add address field
	ID_Card  *multipart.FileHeader `json:"id_card" form:"id_card" binding:"required"` // Add ID Card field
	Username string                `json:"username" form:"username" binding:"required"`
	Password string                `json:"password" form:"password" binding:"required"`
}

type UserUpdateDto struct {
	ID      uuid.UUID             `gorm:"type:char(36);primary_key" json:"id" form:"id"`
	Name    string                `json:"name" form:"name" binding:"required"`
	Email   string                `json:"email" form:"email" binding:"required,email"`
	NoTelp  string                `json:"no_telp" form:"no_telp" binding:"required"`
	Address string                `json:"address" form:"address" binding:"required"` // Add address field
	ID_Card *multipart.FileHeader `json:"id_card" form:"id_card" binding:"required"` // Add ID Card field
}

type UserLoginDTO struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type UserRequestDataDTO struct {
	Username          string `json:"username" form:"username"`
	EncyptedSecretKey string `json:"encrypted_secret_key"`
}

type UserSymmetricKeysResponseDTO struct {
	SecretKey      string `json:"secret_key"`
	SecretKey8Byte string `json:"secret_key_8_byte"`
}

type EncryptedRequestedUserSymmetricKeysDTO struct {
	EncryptedSecretKey string `json:"encrypted_secret_key"`
}

type DecryptedRequestedUserSymmetricKeysDTO struct {
	DecryptedSecretKey string `json:"decrypted_secret_key"`
}
