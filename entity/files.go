package entity

import (
	"ki-d-assignment/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Files struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name      string    `json:"name" binding:"required"`
	Files_AES string    `json:"files_aes" binding:"required"`

	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;foreignKey:Users"`
	Timestamp
}

func (f *Files) BeforeCreate(tx *gorm.DB) error {
	if encrypted, err := utils.EncryptAESCBCFile([]byte(f.Files_AES)); err == nil {
		f.Files_AES = string(encrypted)
	}

	return nil
}