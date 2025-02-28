package entity

// UserDomain handles the many-to-many relationship between User and Domain
type UserDomain struct {
	UserID   uint64 `gorm:"primaryKey"`
	DomainID uint64 `gorm:"primaryKey"`

	User   User   `json:"user" gorm:"foreignKey:UserID"`
	Domain Domain `json:"domain" gorm:"foreignKey:DomainID"`
}
