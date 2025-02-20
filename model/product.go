package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint64         `json:"id" gorm:"not null" groups:"admin"`
	Title       string         `json:"title" gorm:"not null;size:255" groups:"admin"`
	Description string         `json:"description" gorm:"not null;type:text;" groups:"admin"`
	IsActive    bool           `json:"is_active" gorm:"not null" groups:"admin"`
	IconPath    string         `json:"icon_path" gorm:"not null;size:255" groups:"admin"`
	IconUrl     string         `json:"icon_url" gorm:"not null;size:255" groups:"admin"`
	CreatedAt   time.Time      `json:"created_at" gorm:"not null" groups:"admin"`
	UpdatedAt   time.Time      `json:"-" gorm:"not null"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}
