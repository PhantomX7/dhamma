package seed

import (
	"os"

	"github.com/PhantomX7/dhamma/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedRootUser(db *gorm.DB) error {
	users := []entity.User{
		{
			Username: os.Getenv("ADMIN_USERNAME"),
			IsActive: true,
		},
	}

	for _, user := range users {
		if db.First(&entity.User{}, entity.User{
			Username: user.Username,
		}).Error != gorm.ErrRecordNotFound {
			continue
		}

		password, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PASSWORD")), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(password)

		err = db.Create(&user).Error
		if err != nil {
			return err
		}
	}

	return nil
}
