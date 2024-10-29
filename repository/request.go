package repository

import (
	"context"
	"ki-d-assignment/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccessRequestRepository interface {
	CreateAccessRequest(ctx context.Context, request entity.AccessRequest) (entity.AccessRequest, error)
	GetAccessRequestByUsers(ctx context.Context, userID, allowedUserID uuid.UUID) (entity.AccessRequest, error)
}

type accessRequestConnection struct {
	connection *gorm.DB
}

func NewAccessRequestRepository(db *gorm.DB) AccessRequestRepository {
	return &accessRequestConnection{
		connection: db,
	}
}

func (db *accessRequestConnection) CreateAccessRequest(ctx context.Context, request entity.AccessRequest) (entity.AccessRequest, error) {
	request.ID = uuid.New()
	err := db.connection.Create(&request).Error
	return request, err
}

func (db *accessRequestConnection) GetAccessRequestByUsers(ctx context.Context, userID, allowedUserID uuid.UUID) (entity.AccessRequest, error) {
	var request entity.AccessRequest
	err := db.connection.Where("user_id = ? AND allowed_user_id = ?", userID, allowedUserID).First(&request).Error
	return request, err
}
