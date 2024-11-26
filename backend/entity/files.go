package entity

import (
	"ki-d-assignment/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Files struct {
	ID        uuid.UUID `gorm:"primary_key;not_null;type:char(36)" json:"id"`
	Files_AES string    `json:"files_aes" binding:"required"`
	Files_RC4 string    `json:"files_rc4" binding:"required"`
	Files_DES string    `json:"files_des" binding:"required"`
	Signature string    `gorm:"type:text;" json:"signature"`
	// take the user_id from users table
	UserID uuid.UUID `gorm:"type:char(36);not_null;" json:"user_id"`
	Timestamp
}

func (f *Files) BeforeCreate(tx *gorm.DB) error {
	var err error

	// Get the user's encryption keys (from the User associated with this file)
	user := User{}
	if err := tx.Where("id = ?", f.UserID).First(&user).Error; err != nil {
		return err
	}

	// Encrypt the file paths using AES, RC4, and DES
	f.Files_AES, _, _, err = helpers.EncryptData(f.Files_AES, user.SecretKey, user.SecretKey8Byte)
	if err != nil {
		return err
	}

	_, f.Files_DES, _, err = helpers.EncryptData(f.Files_DES, user.SecretKey, user.SecretKey8Byte)
	if err != nil {
		return err
	}

	_, _, f.Files_RC4, err = helpers.EncryptData(f.Files_RC4, user.SecretKey, user.SecretKey8Byte)
	if err != nil {
		return err
	}

	return nil
}

func (f *Files) BeforeUpdate(tx *gorm.DB) error {
	// Re-encrypt the file paths during updates
	return f.BeforeCreate(tx)
}
