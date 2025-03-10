package entity

import (
	"time"

	"github.com/golang-jwt/jwt/v4"

	"gorm.io/gorm"
)

type Timestamp struct {
	CreatedAt time.Time      `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"not null"`
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
