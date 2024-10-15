package entity

import (
	"github.com/google/uuid"
)

type Files struct {
	ID        uuid.UUID `gorm:"primary_key;not_null;type:char(36)" json:"id"`
	Name      string    `json:"name" binding:"required"`
	Files_AES string    `json:"files_aes" binding:"required"`
	SecretKey string    `json:"secret" binding:"required"`
	IV        string    `json:"iv" binding:"required"`
	Signature string    `json:"signature" binding:"required"`
	// take the user_id from users table
	UserID uuid.UUID `gorm:"foreignKey;type:char(36)" json:"user_id"`
	Timestamp
}
