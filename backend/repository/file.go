package repository

import (
	"context"
	"ki-d-assignment/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileRepository interface {
	UploadFile(ctx context.Context, file entity.Files) (entity.Files, error)
	GetFilesByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Files, error)
	GetFileByIDAndUserID(ctx context.Context, fileID uuid.UUID, userID uuid.UUID) (entity.Files, error)
	UpdateFile(file entity.Files) error
}

type fileConnection struct {
	connection *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileConnection{
		connection: db,
	}
}

func (db *fileConnection) UploadFile(ctx context.Context, file entity.Files) (entity.Files, error) {
	file.ID = uuid.New()
	fc := db.connection.Create(&file)
	if fc.Error != nil {
		return entity.Files{}, fc.Error
	}
	return file, nil
}

func (db *fileConnection) GetFilesByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Files, error) {
	var files []entity.Files
	fc := db.connection.Where("user_id = ? AND deleted_at IS NULL", userID).Find(&files)
	if fc.Error != nil {
		return nil, fc.Error
	}
	return files, nil
}

func (db *fileConnection) GetFileByIDAndUserID(ctx context.Context, fileID uuid.UUID, userID uuid.UUID) (entity.Files, error) {
	var file entity.Files
	fc := db.connection.Where("id = ? AND user_id = ? AND deleted_at IS NULL", fileID, userID).Take(&file)
	if fc.Error != nil {
		return entity.Files{}, fc.Error
	}
	return file, nil
}

func (db *fileConnection) UpdateFile(file entity.Files) error {
	if err := db.connection.Save(&file).Error; err != nil {
		return err
	}
	return nil
}
