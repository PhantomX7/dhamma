package transaction_manager

import "gorm.io/gorm"

type Client interface {
	NewTransaction() *gorm.DB
	ExecuteInTransaction(fn TransactionFunc) error
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

type TransactionFunc func(tx *gorm.DB) error

func (tm *client) ExecuteInTransaction(fn TransactionFunc) error {
	tx := tm.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
