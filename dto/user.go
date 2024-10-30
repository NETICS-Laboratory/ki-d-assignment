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
type UserLoginDTO struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type UserRequestDecryptedDto struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	NoTelp   string    `json:"no_telp"`
	Address  string    `json:"address"`
	ID_Card  string    `json:"id_card"`
}

type AccessRequestCreateDto struct {
	RequestedUsername string `json:"requested_username" binding:"required"`
}

type AccessRequestResponseDto struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	RequestedUserID uuid.UUID `json:"requested_user_id"`
	Status          string    `json:"status"`
}

type AccessRequestChangeStatusDto struct {
	Status string `json:"status" form:"status" binding:"required" validate:"oneof=pending approved denied"`
}

type RSAKeyDecryptDto struct {
	EncryptedKey      string `json:"encrypted_key" binding:"required"`
	EncryptedKey8Byte string `json:"encrypted_key_8_byte" binding:"required"`
}

type RSAKeyDecryptResponseDto struct {
	DecryptedKey      string `json:"decrypted_key"`
	DecryptedKey8Byte string `json:"decrypted_key_8_byte"`
}

type AccessPrivateDataRequestDto struct {
	RequestedUserUsername string `json:"requested_user_username" binding:"required"`
	SecretKey             string `json:"secret_key" binding:"required"`
	SecretKey8Byte        string `json:"secret_key_8_byte" binding:"required"`
}
