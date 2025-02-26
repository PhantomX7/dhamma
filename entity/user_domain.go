package entity

// UserDomain handles the many-to-many relationship between User and Domain
type UserDomain struct {
	UserID   uint `gorm:"primaryKey"`
	DomainID uint `gorm:"primaryKey"`

	User   User   `json:"user" gorm:"foreignKey:UserID"`
	Domain Domain `json:"domain" gorm:"foreignKey:DomainID"`
}
