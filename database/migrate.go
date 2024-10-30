package database

import (
	"fmt"
	"ki-d-assignment/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		entity.User{},
		entity.Files{},
		entity.AccessRequest{},
	); err != nil {
		return err
	}
	fmt.Println("Migration success!")

	return nil
}
