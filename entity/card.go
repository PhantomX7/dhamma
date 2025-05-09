package entity

// Card represents a card associated with a follower.
// Each card has a unique code.
type Card struct {
	ID         uint64 `json:"id" gorm:"primary_key;not null"`
	DomainID   uint64 `json:"domain_id" gorm:"not null;index"`      // Foreign key to DomainInf
	FollowerID uint64 `json:"follower_id" gorm:"not null;index"`    // Foreign key to Follower
	Code       string `json:"code" gorm:"not null;unique;size:100"` // Unique identifier for the card
	Timestamp

	// Follower is the follower to whom this card belongs.
	Follower *Follower `json:"follower,omitempty" gorm:"foreignKey:FollowerID"`
	Domain   *Domain   `json:"domain,omitempty" gorm:"foreignKey:DomainID"`
}

// TableName specifies the table name for the Card entity.
func (Card) TableName() string {
	return "cards"
}
