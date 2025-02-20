package model

import (
	"time"

	"gorm.io/gorm"
)

type PostImage struct {
	ID        uint64         `json:"id" gorm:"not null" groups:"admin"`
	PostID    uint64         `json:"post_id" gorm:"not null" groups:"admin"`
	ImagePath string         `json:"image_path" gorm:"not null;size:255" groups:"admin"`
	ImageUrl  string         `json:"image_url" gorm:"not null;size:255" groups:"admin"`
	CreatedAt time.Time      `json:"created_at" gorm:"not null" groups:"admin"`
	UpdatedAt time.Time      `json:"-" gorm:"not null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
