package dto

import (
	"github.com/google/uuid"
	"mime/multipart"
)

type UserCreateDto struct {
	Username 	string				`json:"username" binding:"required"`
	Password 	string      		`json:"password" binding:"required"`
}

type UserUpdateDto struct {
	ID       	uuid.UUID             `gorm:"primary_key" json:"id"`
	Name     	string                `json:"name" binding:"required"`
	Address   	string                `json:"number" binding:"required"`
	CV       	*multipart.FileHeader `json:"cv" binding:"required"`
	ID_Card  	*multipart.FileHeader `json:"id_card" binding:"required"`
	Video    	*multipart.FileHeader `json:"video" binding:"required"`
}

type UserLoginDTO struct {
	ID       	uuid.UUID             `gorm:"primary_key" json:"id"`
	Username 	string                `json:"username" binding:"required"`
	Password 	string                `json:"password" binding:"required"`
}