package entity

import (
	"time"

	"gorm.io/gorm"
)

type Timestamp struct {
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

type Model struct {
	ID uint64 `json:"id" gorm:"primary_key;not null"`

	Timestamp
}
