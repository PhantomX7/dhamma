package entity

// UserRole handles the many-to-many relationship between User and Role
// with domain context
type UserRole struct {
	Model
	UserID   uint64 `gorm:"index:idx_user_domain_role"`
	DomainID uint64 `gorm:"index:idx_user_domain_role"`
	RoleID   uint64 `gorm:"index:idx_user_domain_role"`

	User   User   `gorm:"foreignKey:UserID"`
	Domain Domain `gorm:"foreignKey:DomainID"`
	Role   Role   `gorm:"foreignKey:RoleID"`
}
