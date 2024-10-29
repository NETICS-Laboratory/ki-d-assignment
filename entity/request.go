package entity

import (
	"github.com/google/uuid"
)

type AccessRequest struct {
	ID            uuid.UUID `gorm:"primary_key;not_null;type:char(36)" json:"id"`
	UserID        uuid.UUID `gorm:"type:char(36)" json:"user_id"`
	AllowedUserID uuid.UUID `gorm:"type:char(36)" json:"allowed_id"`
	Status        string    `json:"status"` // "pending", "approved", "denied"

	Timestamp
}
