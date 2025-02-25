package entity

type Role struct {
	Model
	Name        string `gorm:"size:100;unique;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}
