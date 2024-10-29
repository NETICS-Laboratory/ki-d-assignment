package repository

import (
	"context"
	"ki-d-assignment/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, user entity.User) (entity.User, error)
	// GetAllUser(ctx context.Context) ([]entity.User, error)
	FindUserByUsername(ctx context.Context, username string) (entity.User, error)
	FindUserByID(ctx context.Context, userID uuid.UUID) (entity.User, error)
	// DeleteUser(ctx context.Context, userID uuid.UUID) (error)
	// UpdateUser(ctx context.Context, user entity.User) (error)
	CreateAccessRequest(ctx context.Context, accessRequest entity.AccessRequest) (entity.AccessRequest, error)
	GetAccessRequestsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.AccessRequest, error)
	UpdateAccessRequestStatus(ctx context.Context, requestID uuid.UUID, status string) error
}

type userConnection struct {
	connection *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) RegisterUser(ctx context.Context, user entity.User) (entity.User, error) {
	uc := db.connection.Create(&user)
	if uc.Error != nil {
		return entity.User{}, uc.Error
	}
	return user, nil
}

// func(db *userConnection) GetAllUser(ctx context.Context) ([]entity.User, error) {
// 	var listUser []entity.User
// 	tx := db.connection.Find(&listUser)
// 	if tx.Error != nil {
// 		return nil, tx.Error
// 	}
// 	return listUser, nil
// }

func (db *userConnection) FindUserByUsername(ctx context.Context, username string) (entity.User, error) {
	var user entity.User
	ux := db.connection.Where("username = ?", username).Take(&user)
	if ux.Error != nil {
		return user, ux.Error
	}
	return user, nil
}

func (db *userConnection) FindUserByID(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	var user entity.User
	ux := db.connection.Where("id = ?", userID).Take(&user)
	if ux.Error != nil {
		return user, ux.Error
	}
	return user, nil
}

func (db *userConnection) CreateAccessRequest(ctx context.Context, accessRequest entity.AccessRequest) (entity.AccessRequest, error) {
	accessRequest.ID = uuid.New()
	result := db.connection.Create(&accessRequest)
	if result.Error != nil {
		return entity.AccessRequest{}, result.Error
	}
	return accessRequest, nil
}

func (db *userConnection) GetAccessRequestsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.AccessRequest, error) {
	var requests []entity.AccessRequest
	err := db.connection.Where("allowed_user_id = ?", userID).Find(&requests).Error
	return requests, err
}

func (db *userConnection) UpdateAccessRequestStatus(ctx context.Context, requestID uuid.UUID, status string) error {
	return db.connection.Model(&entity.AccessRequest{}).Where("id = ?", requestID).Update("status", status).Error
}

// func(db *userConnection) DeleteUser(ctx context.Context, userID uuid.UUID) (error) {
// 	uc := db.connection.Delete(&entity.User{}, &userID)
// 	if uc.Error != nil {
// 		return uc.Error
// 	}
// 	return nil
// }

// func(db *userConnection) UpdateUser(ctx context.Context, user entity.User) (error) {
// 	uc := db.connection.Updates(&user)
// 	if uc.Error != nil {
// 		return uc.Error
// 	}
// 	return nil
// }
