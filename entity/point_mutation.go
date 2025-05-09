package entity

// Constants for PointMutation SourceType
// For GORM polymorphic associations, SourceType should typically store the table name of the source entity.
const (
	// PointMutationSourceTypeEventAttendance indicates the point mutation originated from an event attendance.
	// This should match the table name of the EventAttendance entity.
	PointMutationSourceTypeEventAttendance = "event_attendances"
	// PointMutationSourceTypeManual indicates a manual adjustment by an admin or system process.
	PointMutationSourceTypeManual = "manual_adjustment"
	// PointMutationSourceTypeInitial indicates points assigned at follower creation or initial setup.
	PointMutationSourceTypeInitial = "initial_setup"
	// Add other source types as needed
)

// PointMutation tracks changes to a follower's points.
// This entity uses a polymorphic association for its source, meaning a point mutation
// can originate from different types of entities (e.g., EventAttendance) or system processes.
type PointMutation struct {
	ID          uint64  `json:"id" gorm:"primary_key;not null"`
	FollowerID  uint64  `json:"follower_id" gorm:"not null;index"` // Foreign key to Follower
	Amount      int     `json:"amount" gorm:"not null"`            // Can be positive or negative
	SourceType  string  `json:"source_type" gorm:"not null;size:100;index"`
	SourceID    *uint64 `json:"source_id,omitempty" gorm:"null;index"` // Optional ID of the source (e.g., EventAttendanceID)
	Description *string `json:"description,omitempty" gorm:"size:255;null"`
	Timestamp

	// Follower is the follower to whom these points belong.
	Follower *Follower `json:"follower,omitempty" gorm:"foreignKey:FollowerID"`
}

// TableName specifies the table name for the PointMutation entity.
func (PointMutation) TableName() string {
	return "point_mutations"
}
