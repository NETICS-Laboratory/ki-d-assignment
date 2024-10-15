package dto

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type FileCreateDto struct {
	ID   uuid.UUID             `gorm:"primary_key;not_null;type:char(36)" json:"id"`
	File *multipart.FileHeader `json:"file" form:"file" binding:"required"`
}
