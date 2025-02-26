package entity

type Domain struct {
	Model
	Name        string `gorm:"size:100;unique;not null" json:"name"`
	Code        string `gorm:"size:50;unique;not null" json:"code"`
	Description string `gorm:"size:255" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`

	// Has-Many relationship with Role
	Roles []Role `gorm:"foreignKey:DomainID"`
	// Many-to-Many with User through UserDomain
	Users []User `gorm:"many2many:user_domains;"`
}
