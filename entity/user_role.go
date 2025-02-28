package entity

// UserRole handles the many-to-many relationship between User and Role
// with domain context
type UserRole struct {
	ID       uint64 `json:"id" gorm:"primary_key;not null"`
	UserID   uint64 `gorm:"index:idx_user_domain_role" json:"user_id"`
	DomainID uint64 `gorm:"index:idx_user_domain_role" json:"domain_id"`
	RoleID   uint64 `gorm:"index:idx_user_domain_role" json:"role_id"`
	Timestamp

	User   User   `gorm:"foreignKey:UserID" json:"user"`
	Domain Domain `gorm:"foreignKey:DomainID" json:"domain"`
	Role   Role   `gorm:"foreignKey:RoleID" json:"role"`
}
