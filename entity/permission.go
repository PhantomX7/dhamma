package entity

type Permission struct {
	ID               uint64  `json:"id" gorm:"primary_key;not null"`
	Name             string  `json:"name" gorm:"size:100;unique;not null"`
	Code             string  `json:"code" gorm:"size:100;unique;not null"`
	Object           string  `json:"object" gorm:"size:255;not null"`
	Action           string  `json:"action" gorm:"size:50;not null"`
	Type             string  `json:"type" gorm:"size:50;not null"`
	Description      string  `json:"description" gorm:"size:255"`
	IsDomainSpecific bool    `json:"is_domain_specific" gorm:""`
	DomainID         *uint64 `json:"domain_id" gorm:""`
	Timestamp
}
