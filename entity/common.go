package entity

import (
	"github.com/golang-jwt/jwt/v4"
	"time"

	"gorm.io/gorm"
)

type Timestamp struct {
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

type AccessClaims struct {
	UserID uint64 `json:"user_id"`
	Role   string `json:"role"`

	jwt.RegisteredClaims
}

type RefreshClaims struct {
	RefreshToken string `json:"refresh_token"`

	jwt.RegisteredClaims
}
