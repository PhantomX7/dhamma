package entity

type Role struct {
	ID          uint64 `json:"id" gorm:"primary_key;not null"`
	Name        string `gorm:"size:255;index:idx_domain_role_name,unique,priority:2"`
	Description string `gorm:"size:255" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	DomainID    uint64 `gorm:"index:idx_domain_role_name,unique,priority:1"`
	Timestamp
	
	Domain Domain `gorm:"foreignKey:DomainID"`
}
