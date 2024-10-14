package repository

import (
	"context"
	"ki-d-assignment/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, user entity.User) (entity.User, error)
	//GetAllUser(ctx context.Context) ([]entity.User, error)
	// FindUserByEmail(ctx context.Context, email string) (entity.User, error)
	FindUserByID(ctx context.Context, userID uuid.UUID) (entity.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) (error)
	UpdateUser(ctx context.Context, user entity.User) (entity.User, error)
}

type userConnection struct {
	connection *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (db *UserConnection) RegisterUser(ctx context.Context, user entity.User) (entity.User, error) {
	if err := db.connection.Create(&user).Error; err != nil {
		return entity.User{}, err
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

// func(db *userConnection) FindUserByEmail(ctx context.Context, email string) (entity.User, error) {
// 	var user entity.User
// 	ux := db.connection.Where("email = ?", email).Take(&user)
// 	if ux.Error != nil {
// 		return user, ux.Error
// 	}
// 	return user, nil
// }

func(db *userConnection) FindUserByID(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	var user entity.User
	ux := db.connection.Where("id = ?", userID).Take(&user)
	if ux.Error != nil {
		return user, ux.Error
	}
	return user, nil
}

func (db *UserConnection) DeleteUser(ctx context.Context, userID uuid.UUID) (error) {
	if err := db.connection.Where("id = ?", userID).Delete(&entity.User{}).Error; err != nil {
		return err
	}
	return nil
}

func (db *UserConnection) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	if err := db.connection.Save(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}