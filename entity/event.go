package entity

// Event represents an event that followers can attend.
type Event struct {
	ID            uint64  `json:"id" gorm:"primary_key;not null"`
	DomainID      uint64  `json:"domain_id" gorm:"not null;index"` // Foreign key to Domain
	Name          string  `json:"name" gorm:"not null;size:255"`
	Description   *string `json:"description,omitempty" gorm:"type:text;null"`
	PointsAwarded int     `json:"points_awarded" gorm:"not null;default:0"` // Points awarded for attendance
	Timestamp

	Domain           *Domain           `json:"domain,omitempty" gorm:"foreignKey:DomainID"`
	EventAttendances []EventAttendance `json:"event_attendances,omitempty" gorm:"foreignKey:EventID"`
}

// TableName specifies the table name for the Event entity.
func (Event) TableName() string {
	return "events"
}
