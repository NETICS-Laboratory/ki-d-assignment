package dto

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type FileCreateDto struct {
	ID   uuid.UUID             `gorm:"primary_key;not_null;type:char(36)" json:"id"`
	File *multipart.FileHeader `json:"file" form:"file" binding:"required"`
}

type FileDecryptByIDDto struct {
	ID uuid.UUID `json:"id" form:"id" binding:"required"`
}

type FileDecryptedResponse struct {
	ID            uuid.UUID `json:"id"`
	Decrypted_AES string    `json:"decrypted_aes"`
	Decrypted_RC4 string    `json:"decrypted_rc4"`
	Decrypted_DES string    `json:"decrypted_des"`
	Signature     string    `json:"signature"`
}

type FileSignDto struct {
	FileID uuid.UUID `json:"file_id" form:"file_id" binding:"required"`
}

type FileSignResponse struct {
	FileID    uuid.UUID `json:"file_id"`
	Signature string    `json:"signature"`
}

type FileVerifySignatureDto struct {
	FileID    uuid.UUID `json:"file_id" form:"file_id" binding:"required"`
	Signature string    `json:"signature" form:"signature" binding:"required"`
}
