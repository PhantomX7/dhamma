package entity

// UserDomain handles the many-to-many relationship between User and Domain
type UserDomain struct {
	UserID   uint64 `json:"user_id" gorm:"primaryKey"`
	DomainID uint64 `json:"domain_id" gorm:"primaryKey"`

	User   User   `json:"user" gorm:"foreignKey:UserID"`
	Domain Domain `json:"domain" gorm:"foreignKey:DomainID"`
}
