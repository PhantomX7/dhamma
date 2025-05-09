package entity

type Domain struct {
	ID          uint64 `json:"id" gorm:"primary_key;not null"`
	Name        string `json:"name" gorm:"size:100;unique;not null"`
	Code        string `json:"code" gorm:"size:50;unique;not null"`
	Description string `json:"description" gorm:"size:255"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`
	Timestamp

	// Has-Many relationship with Role
	Roles []Role `json:"roles" gorm:"foreignKey:DomainID"`
	// Has-Many relationship with Follower
	Followers []Follower `json:"followers" gorm:"foreignKey:DomainID"`
	// Has-Many relationship with Follower
	Cards []Card `json:"cards" gorm:"foreignKey:DomainID"`
	// Many-to-Many with User through UserDomain
	Users []User `json:"users" gorm:"many2many:user_domains;"`
}
