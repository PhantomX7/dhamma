package entity

// Follower represents a follower within a domain.
type Follower struct {
	ID           uint64  `json:"id" gorm:"primary_key;not null"`
	DomainID     uint64  `json:"domain_id" gorm:"not null;index"` // Foreign key to Domain
	Name         string  `json:"name" gorm:"not null;size:255"`
	Phone        *string `json:"phone" gorm:"size:50;null"` // Optional phone number
	Points       int     `json:"points" gorm:"not null;default:0"`
	IsBloodDonor bool    `json:"is_blood_donor" gorm:"not null;default:false"`
	IsYouth      bool    `json:"is_youth" gorm:"not null;default:false"`
	Timestamp

	Domain           *Domain           `json:"domain,omitempty" gorm:"foreignKey:DomainID"`
	EventAttendances []EventAttendance `json:"event_attendances,omitempty" gorm:"foreignKey:FollowerID"`
	PointMutations   []PointMutation   `json:"point_mutations,omitempty" gorm:"foreignKey:FollowerID"`
	Cards            []Card            `json:"cards,omitempty" gorm:"foreignKey:FollowerID"`
}

// TableName specifies the table name for the Follower entity.
func (Follower) TableName() string {
	return "followers"
}
