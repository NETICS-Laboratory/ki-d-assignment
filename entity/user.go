package entity

import (
	"ki-d-assignment/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	Identity struct {
        name_AES            string      `json:"name_aes" binding:"required"`
        name_DES            string      `json:"name_des" binding:"required"`
        name_RC4            string      `json:"name_rc4" binding:"required"`
        address_AES         string      `json:"address" binding:"required"`
        address_DES         string      `json:"address" binding:"required"`
        address_RC4         string      `json:"address" binding:"required"`
        CV_AES              string      `json:"cv_aes" binding:"required"`
        CV_DES              string      `json:"cv_des" binding:"required"`
        CV_RC4              string      `json:"cv_rc4" binding:"required"`
        ID_card_AES         string      `json:"id_card_aes" binding:"required"`
        ID_card_DES         string      `json:"id_card_des" binding:"required"`
        ID_card_RC4         string      `json:"id_card_rc4" binding:"required"`
		Video_AES           string      `json:"video_aes" binding:"required"`
		Video_DES           string      `json:"video_des" binding:"required"`
		Video_RC4           string      `json:"video_rc4" binding:"required"`
    }

	Credential struct {
        username_AES        string      `json:"username_aes" binding:"required"`
        username_DES        string      `json:"username_des" binding:"required"`
        username_RC4        string      `json:"username_rc4" binding:"required"`
        password_AES        string      `json:"password_aes" binding:"required"`
        password_DES        string      `json:"password_des" binding:"required"`
        password_RC4        string      `json:"password_rc4" binding:"required"`
	}

	Key struct {
        secret_key          string      `json:"secret_key" binding:"required"`
        IV                  string      `json:"iv" binding:"required"`
        secret_key_8_byte   string      `json:"secret_key_8_byte" binding:"required"`
        IV_8_byte           string      `json:"iv_8_byte" binding:"required"`
    }

)

type User struct {
    ID  					uuid.UUID 	`gorm:"primary_key;not_null;type:char(36)" json:"id"`
    Identity
    Credential
    Key

    Files           []Files             `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" binding:"required" json:"files"`
    AllowedUsers    []AllowedUser       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" binding:"required" json:"allowed_users"`

    Timestamp
}

func (User) TableName() string {
	return "users"
}

func (u *user) EncryptUsername (username string) string {
	enc, - := utils.EncryptAES([]byte(username), []byte(u.Username_AES))
	return string(enc)
}

func (u *User) RC4EncryptField(field *string) error {
    enc, err := utils.EncryptRC4([]byte(*field), u.secret_key)
    if err != nil {
        return err
    }
    *field = string(enc)
    return nil
}

func (u *User) EncryptAllFieldsRC4() error {
    fields := []*string{
        &u.username_RC4,
        &u.password_RC4,
        &u.name_RC4,
        &u.address_RC4,
        &u.CV_RC4,
        &u.ID_Card_RC4,
    }

    for _, field := range fields {
        if err := u.RC4EncryptField(field); err != nil {
            return err
        }
    }
    
    return nil 
}


