package transaction_manager

import "gorm.io/gorm"

type Client interface {
	NewTransaction() *gorm.DB
}

type client struct {
	db *gorm.DB
}

func New(db *gorm.DB) Client {
	return &client{
		db: db,
	}
}

func (tm *client) NewTransaction() *gorm.DB {
	return tm.db.Begin()
}
