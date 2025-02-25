package entity

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `json:"id" gorm:"primary_key;not null"`
	UserID    uint64    `json:"user_id" gorm:"not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	IsValid   bool      `json:"is_valid" gorm:"not null;default:true"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	User User `gorm:"foreignKey:UserID"`
}
