package entity

type Role struct {
	ID          uint64 `json:"id" gorm:"primary_key;not null"`
	DomainID    uint64 `json:"domain_id" gorm:"index:idx_domain_role_name,unique,priority:1"`
	Name        string `json:"name" gorm:"size:255;index:idx_domain_role_name,unique,priority:2"`
	Description string `json:"description" gorm:"size:255"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`
	Timestamp

	Domain *Domain `json:"domain,omitempty" gorm:"foreignKey:DomainID"`
}
