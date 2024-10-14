package entity

import (
	"gorm.io/gorm"
	"github.com/google/uuid"
)
type Files struct {
	gorm.Model

	Files 		[]byte 			`json:"files" binding:"required"`
	UserID 		uuid.UUID 		`gorm:"foreignKey;type:char(36)" json:"user_id"`
}