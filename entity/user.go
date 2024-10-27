package entity

import (
	"ki-d-assignment/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	Identity struct {
		Name_AES    string    `json:"name_aes" binding:"required"`
		Name_DES    string    `json:"name_des" binding:"required"`
		Name_RC4    string    `json:"name_rc4" binding:"required"`
		Email_AES   string    `json:"email_aes" binding:"required"`
		Email_DES   string    `json:"email_des" binding:"required"`
		Email_RC4   string    `json:"email_rc4" binding:"required"`
		NoTelp_AES  string    `json:"notelp_aes" binding:"required"`
		NoTelp_DES  string    `json:"notelp_des" binding:"required"`
		NoTelp_RC4  string    `json:"notelp_rc4" binding:"required"`
		Address_AES string    `json:"address_aes" binding:"required"`
		Address_DES string    `json:"address_des" binding:"required"`
		Address_RC4 string    `json:"address_rc4" binding:"required"`
		ID_Card_ID  uuid.UUID `json:"id_card_id"`
		ID_Card_AES string    `json:"id_card_aes" binding:"required"`
		ID_Card_DES string    `json:"id_card_des" binding:"required"`
		ID_Card_RC4 string    `json:"id_card_rc4" binding:"required"`
	}

	Credential struct {
		Username     string `json:"username" binding:"required"`
		Username_AES string `json:"username_aes" binding:"required"`
		Username_DES string `json:"username_des" binding:"required"`
		Username_RC4 string `json:"username_rc4" binding:"required"`
		Password_AES string `json:"password_aes" binding:"required"`
		Password_DES string `json:"password_des" binding:"required"`
		Password_RC4 string `json:"password_rc4" binding:"required"`
	}

	Key struct {
		SecretKey      []byte `json:"secret" binding:"required"`
		SecretKey8Byte []byte `json:"secret_key_8_byte" binding:"required"`
	}
)

type User struct {
	ID uuid.UUID `gorm:"primary_key;not_null;type:char(36)" json:"id"`
	Identity
	Credential
	Key

	Files []Files `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" binding:"required" json:"files"`

	Timestamp
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	var err error

	// Username encryption
	u.Username_AES, u.Username_DES, u.Username_RC4, err = helpers.EncryptData(u.Username, u.SecretKey, u.SecretKey8Byte)
	if err != nil {
		return err
	}

	// Password encryption
	u.Password_AES, u.Password_DES, u.Password_RC4, err = helpers.EncryptData(u.Password_AES, u.SecretKey, u.SecretKey8Byte)
	if err != nil {
		return err
	}

	// Identity encryption (Name, Email, NoTelp, Address, ID Card)
	u.Name_AES, u.Name_DES, u.Name_RC4, err = helpers.EncryptData(u.Name_AES, u.SecretKey, u.SecretKey8Byte)
	if err != nil {
		return err
	}

	u.Email_AES, u.Email_DES, u.Email_RC4, err = helpers.EncryptData(u.Email_AES, u.SecretKey, u.SecretKey8Byte)
	if err != nil {
		return err
	}

	u.NoTelp_AES, u.NoTelp_DES, u.NoTelp_RC4, err = helpers.EncryptData(u.NoTelp_AES, u.SecretKey, u.SecretKey8Byte)
	if err != nil {
		return err
	}

	u.Address_AES, u.Address_DES, u.Address_RC4, err = helpers.EncryptData(u.Address_AES, u.SecretKey, u.SecretKey8Byte)
	if err != nil {
		return err
	}

	//encrypt link to id card
	u.ID_Card_AES, _, _, err = helpers.EncryptData(u.ID_Card_AES, u.SecretKey, u.SecretKey8Byte)
	// fmt.Println(u.ID_Card_AES)
	if err != nil {
		return err
	}

	_, u.ID_Card_DES, _, err = helpers.EncryptData(u.ID_Card_DES, u.SecretKey, u.SecretKey8Byte)
	// fmt.Println(u.ID_Card_DES)
	if err != nil {
		return err
	}

	_, _, u.ID_Card_RC4, err = helpers.EncryptData(u.ID_Card_RC4, u.SecretKey, u.SecretKey8Byte)
	// fmt.Println(u.ID_Card_RC4)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// Same encryption logic for update as in BeforeCreate
	return u.BeforeCreate(tx)
}
