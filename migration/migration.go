package migration

import (
	"github.com/PhantomX7/dhamma/entity"

	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) error {
	return db.AutoMigrate(
		// list all migration here
		entity.User{},
	)
}
