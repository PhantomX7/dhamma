package seed

import (
	"errors"
	"fmt"
	"log"

	"github.com/PhantomX7/dhamma/constants"
	"github.com/PhantomX7/dhamma/entity"
	"gorm.io/gorm"
)

type PermissionSeeder struct {
	db *gorm.DB
}

func NewPermissionSeeder(db *gorm.DB) *PermissionSeeder {
	return &PermissionSeeder{db: db}
}

func (s *PermissionSeeder) GenerateApiPermissions() (err error) {
	log.Print("seeding api permissions")
	for _, permission := range constants.ApiPermissions {
		permissionCode := fmt.Sprintf("%s/%s", permission.Object, permission.Action)

		if !errors.Is(s.db.First(&entity.Permission{}, entity.Permission{
			Code: permissionCode,
		}).Error, gorm.ErrRecordNotFound) {
			continue
		}

		permission.Code = permissionCode

		err = s.db.Create(&permission).Error
		if err != nil {
			return err
		}
	}

	return err
}
