package model

import (
	"time"
)

const (
	ConfigKeyEmailDestination = "EMAIL_DESTINATION"
)

type Config struct {
	ID        uint64    `json:"id" gorm:"not null;index" groups:"admin"`
	Key       string    `json:"key" gorm:"not null;size:255;" groups:"admin"`
	Value     string    `json:"value" gorm:"not null;size:255;" groups:"admin"`
	CreatedAt time.Time `json:"created_at" gorm:"not null" groups:"admin"`
	UpdatedAt time.Time `json:"-" gorm:"not null" groups:"admin"`
}
