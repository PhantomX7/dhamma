package entity

import (
	"time"
)

// EventAttendance represents a follower's attendance at an event.
// This acts as a join table between Follower and Event.
type EventAttendance struct {
	ID         uint64    `json:"id" gorm:"primary_key;not null"`
	FollowerID uint64    `json:"follower_id" gorm:"not null;index"`
	EventID    uint64    `json:"event_id" gorm:"not null;index"`
	AttendedAt time.Time `json:"attended_at" gorm:"not null"` // Timestamp of when the follower attended
	Timestamp

	Follower      *Follower      `json:"follower,omitempty" gorm:"foreignKey:FollowerID"`
	Event         *Event         `json:"event,omitempty" gorm:"foreignKey:EventID"`
	PointMutation *PointMutation `json:"point_mutation,omitempty" gorm:"polymorphic:Source;"` // "Source" matches the prefix of SourceID & SourceType
}

// TableName specifies the table name for the EventAttendance entity.
func (EventAttendance) TableName() string {
	return "event_attendances"
}
