package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uint64    `json:"id" gorm:"primary_key;not null"`
	UUID     uuid.UUID `json:"uuid" gorm:"not null"`
	Username string    `json:"username" gorm:"not null;size:255"`
	Password string    `json:"-" gorm:"not null;size:255"`
	IsActive bool      `json:"is_active" gorm:"not null"`

	Timestamp
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	user.UUID = uuid.New()
	return
}
