package database_transaction

import "gorm.io/gorm"

type TransactionManager interface {
	NewTransaction() *gorm.DB
}

type transactionManager struct {
	db *gorm.DB
}

func New(db *gorm.DB) TransactionManager {
	return &transactionManager{
		db: db,
	}
}

func (tm *transactionManager) NewTransaction() *gorm.DB {
	return tm.db.Begin()
}
