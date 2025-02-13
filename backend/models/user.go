package models

import "time"

type User struct {
	ID        string    `gorm:"uniqueIndex;primarykey" json:"id"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserAPIKey struct {
	UserID          string `gorm:"uniqueIndex;primarykey" json:"user_id"`
	EncryptedAPIKey string `gorm:"column:encrypted_api_key" json:"-"`
}
