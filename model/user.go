package model

import (
	"time"
)

type User struct {
	ID        uint64    `json:"id" gorm:"not null" groups:"admin"`
	Username  string    `json:"username" gorm:"not null;size:255" groups:"admin"`
	Password  string    `json:"-" gorm:"not null;size:255"`
	IsActive  bool      `json:"is_active" gorm:"not null" groups:"admin"`
	CreatedAt time.Time `json:"created_at" gorm:"not null" groups:"admin"`
	UpdatedAt time.Time `json:"-" gorm:"not null"`
}
