package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Model
	UUID     uuid.UUID `json:"uuid" gorm:"not null"`
	Username string    `json:"username" gorm:"not null;size:255"`
	Password string    `json:"-" gorm:"not null;size:255"`
	IsActive bool      `json:"is_active" gorm:"not null"`

	RefreshTokens []RefreshToken `gorm:"foreignKey:UserID"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	user.UUID = uuid.New()
	return
}
