package entity

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	IsAdmin  bool      `json:"is_admin"`
	// Role     string    `json:"role"`

	Timestamp
}
