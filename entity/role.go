package entity

import "time"

type Role struct {
	Model
	Name        string `gorm:"size:100;unique;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}

type UserRole struct {
	UserID    uint64    `gorm:"primary_key" json:"user_id"`
	RoleID    uint64    `gorm:"primary_key" json:"role_id"`
	Domain    string    `gorm:"size:100;primary_key" json:"domain"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
