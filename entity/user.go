package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Model
	UUID     uuid.UUID `json:"uuid" gorm:"not null"`
	Username string    `json:"username" gorm:"uniqueIndex;not null;size:255"`
	Password string    `json:"-" gorm:"not null;size:255"`
	IsActive bool      `json:"is_active" gorm:"not null"`

	// Many-to-Many with Domain through UserDomain
	Domains []Domain `gorm:"many2many:user_domains;"`
	// Has-Many with UserRole to handle multiple roles per domain
	UserRoles []UserRole `gorm:"foreignKey:UserID"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	user.UUID = uuid.New()
	return
}
