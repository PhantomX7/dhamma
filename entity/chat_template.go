package entity

// ChatTemplate represents a chat template that can be used within a domain.
// Each domain can have multiple chat templates with one marked as default.
type ChatTemplate struct {
	ID          uint64  `json:"id" gorm:"primary_key;not null"`
	DomainID    uint64  `json:"domain_id" gorm:"not null;index"` // Foreign key to Domain
	Name        string  `json:"name" gorm:"not null;size:255"`
	Description *string `json:"description" gorm:"size:500;null"`
	Content     string  `json:"content" gorm:"not null;type:text"` // Template content/body
	IsDefault   bool    `json:"is_default" gorm:"not null;default:false"`
	IsActive    bool    `json:"is_active" gorm:"not null;default:true"`
	Timestamp

	Domain *Domain `json:"domain,omitempty" gorm:"foreignKey:DomainID"`
}

// TableName specifies the table name for the ChatTemplate entity.
func (ChatTemplate) TableName() string {
	return "chat_templates"
}
