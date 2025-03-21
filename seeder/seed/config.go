package seed

import (
	"errors"

	"github.com/PhantomX7/dhamma/entity"

	"gorm.io/gorm"
)

func SeedConfig(db *gorm.DB) error {
	configs := []entity.Config{
		{
			Key:   entity.ConfigKeyEmailDestination,
			Value: "",
		},
	}

	for _, config := range configs {
		if !errors.Is(db.First(&entity.Config{}, entity.Config{
			Key:   config.Key,
			Value: config.Value,
		}).Error, gorm.ErrRecordNotFound) {
			continue
		}

		err := db.Create(&config).Error
		if err != nil {
			return err
		}
	}

	return nil
}
