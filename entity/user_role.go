package entity

// UserRole handles the many-to-many relationship between User and Role
// with domain context
type UserRole struct {
	ID       uint64 `json:"id" gorm:"primary_key;not null"`
	UserID   uint64 `json:"user_id" gorm:"index:idx_user_domain_role"`
	DomainID uint64 `json:"domain_id" gorm:"index:idx_user_domain_role"`
	RoleID   uint64 `json:"role_id" gorm:"index:idx_user_domain_role"`
	Timestamp

	User   *User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Domain *Domain `json:"domain,omitempty" gorm:"foreignKey:DomainID"`
	Role   *Role   `json:"role,omitempty" gorm:"foreignKey:RoleID"`
}
