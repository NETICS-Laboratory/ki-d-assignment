package entity

import (
	// "ki-d-assignment/helpers"
	"ki-d-assignment/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	Identity struct {
		Name_AES    string `json:"name_aes" binding:"required"`
		Address_AES string `json:"address_aes" binding:"required"`
		CV_AES      string `json:"cv_aes" binding:"required"`
		ID_Card_AES string `json:"id_card_aes" binding:"required"`
	}

	Credential struct {
		Username     string `json:"username" binding:"required"`
		Username_AES string `json:"username_aes" binding:"required"`
		Password_AES string `json:"password_aes" binding:"required"`
	}
)

type User struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Identity
	Credential

	Files []Files `json:"files" gorm:"foreignKey:UserID" binding:"required"`
}

// TODO: Fix encrypt parameters, Add Decrypt, Testing
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if encrypted, err := utils.EncryptAESCBC(u.Username_AES); err == nil {
		u.Username_AES = string(encrypted)
	}

	if encrypted, err := utils.EncryptAESCBC(u.Password_AES); err == nil {
		u.Password_AES = string(encrypted)
	}

	if encrypted, err := utils.EncryptAESCBC(u.Name_AES); err == nil {
		u.Name_AES = string(encrypted)
	}

	if encrypted, err := utils.EncryptAESCBC(u.Address_AES); err == nil {
		u.Address_AES = string(encrypted)
	}

	if encrypted, err := utils.EncryptAESCBC(u.CV_AES); err == nil {
		u.CV_AES = string(encrypted)
	}

	if encrypted, err := utils.EncryptAESCBC(u.ID_Card_AES); err == nil {
		u.ID_Card_AES = string(encrypted)
	}
	return nil

}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if encrypted, err := utils.EncryptAESCBC(u.Username_AES); err == nil {
		u.Username_AES = string(encrypted)
	}

	if encrypted, err := utils.EncryptAESCBC(u.Password_AES); err == nil {
		u.Password_AES = string(encrypted)
	}

	if encrypted, err := utils.EncryptAESCBC(u.Name_AES); err == nil {
		u.Name_AES = string(encrypted)
	}

	if encrypted, err := utils.EncryptAESCBC(u.Address_AES); err == nil {
		u.Address_AES = string(encrypted)
	}

	if encrypted, err := utils.EncryptAESCBC(u.CV_AES); err == nil {
		u.CV_AES = string(encrypted)
	}

	if encrypted, err := utils.EncryptAESCBC(u.ID_Card_AES); err == nil {
		u.ID_Card_AES = string(encrypted)
	}

	return nil
}
