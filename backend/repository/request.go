package repository

import (
	"context"
	"ki-d-assignment/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccessRequestRepository interface {
	CreateAccessRequest(ctx context.Context, request entity.AccessRequest) (entity.AccessRequest, error)
	GetAccessRequestsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.AccessRequest, error)
	GetAccessRequestsByID(ctx context.Context, requestID uuid.UUID) (entity.AccessRequest, error)
	GetAccessRequestsByRequestedUserID(ctx context.Context, requestedUserID uuid.UUID) ([]entity.AccessRequest, error)
	CheckExistingAccessRequest(ctx context.Context, userID, requestedUserID uuid.UUID) (bool, error)
	UpdateAccessRequestStatus(ctx context.Context, requestID uuid.UUID, status string) error
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

func (db *accessRequestConnection) CheckExistingAccessRequest(ctx context.Context, userID, requestedUserID uuid.UUID) (bool, error) {
	var request entity.AccessRequest
	err := db.connection.Where("user_id = ? AND requested_user_id = ?", userID, requestedUserID).Take(&request).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (db *accessRequestConnection) GetAccessRequestsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.AccessRequest, error) {
	var requests []entity.AccessRequest
	err := db.connection.Where("user_id = ?", userID).Find(&requests).Error
	return requests, err
}

func (db *accessRequestConnection) GetAccessRequestsByRequestedUserID(ctx context.Context, requestedUserID uuid.UUID) ([]entity.AccessRequest, error) {
	var requests []entity.AccessRequest
	err := db.connection.Where("requested_user_id = ?", requestedUserID).Find(&requests).Error
	return requests, err
}

func (db *accessRequestConnection) GetAccessRequestsByID(ctx context.Context, requestID uuid.UUID) (entity.AccessRequest, error) {
	var request entity.AccessRequest
	err := db.connection.Where("id = ?", requestID).Take(&request).Error
	return request, err
}

func (db *accessRequestConnection) UpdateAccessRequestStatus(ctx context.Context, requestID uuid.UUID, status string) error {
	return db.connection.Model(&entity.AccessRequest{}).Where("id = ?", requestID).Update("status", status).Error
}
