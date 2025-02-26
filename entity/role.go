package entity

type Role struct {
	Model
	Name        string `gorm:"size:255;index:idx_domain_role_name,unique,priority:2"`
	Description string `gorm:"size:255" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	DomainID    uint64 `gorm:"index:idx_domain_role_name,unique,priority:1"`

	Domain Domain `gorm:"foreignKey:DomainID"`
}
