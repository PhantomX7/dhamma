package utility

import (
	"github.com/PhantomX7/dhamma/model"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) error {
	return db.AutoMigrate(
		model.User{},
		model.Post{},
		model.PostImage{},
		model.Product{},
		model.Config{},
		model.Partner{},
	)
}
