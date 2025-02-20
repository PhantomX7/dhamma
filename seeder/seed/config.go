package seed

import (
	"github.com/PhantomX7/dhamma/model"

	"gorm.io/gorm"
)

func SeedConfig(db *gorm.DB) error {
	configs := []model.Config{
		{
			Key:   model.ConfigKeyEmailDestination,
			Value: "",
		},
	}

	for _, config := range configs {
		if db.First(&model.Config{}, model.Config{
			Key:   config.Key,
			Value: config.Value,
		}).Error != gorm.ErrRecordNotFound {
			continue
		}

		err := db.Create(&config).Error
		if err != nil {
			return err
		}
	}

	return nil
}
