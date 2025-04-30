package seed

import (
	"fmt"
	"log"

	"github.com/PhantomX7/dhamma/constants/permissions"
	"github.com/PhantomX7/dhamma/entity"
	"gorm.io/gorm"
)

// PermissionSeeder handles seeding permissions to the database
type PermissionSeeder struct {
	db *gorm.DB
}

// NewPermissionSeeder creates a new permission seeder
func NewPermissionSeeder(db *gorm.DB) *PermissionSeeder {
	return &PermissionSeeder{db: db}
}

// GenerateApiPermissions seeds API permissions to the database
// It creates new permissions and updates existing ones
func (s *PermissionSeeder) GenerateApiPermissions() (err error) {
	log.Print("seeding api permissions")

	// Begin transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Defer rollback in case of error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("panic occurred during permission seeding: %v", r)
			}
		} else if err != nil {
			tx.Rollback()
		}
	}()

	// Get existing permissions for bulk checking
	var existingPermissions []entity.Permission
	if err = tx.Find(&existingPermissions).Error; err != nil {
		return err
	}

	// Create map for quick lookup
	existingMap := make(map[string]entity.Permission)
	for _, perm := range existingPermissions {
		existingMap[perm.Code] = perm
	}

	// Process permissions
	for _, permission := range permissions.ApiPermissions {
		permissionCode := fmt.Sprintf("%s/%s", permission.Object, permission.Action)
		permission.Code = permissionCode

		// Check if permission exists
		if existing, found := existingMap[permissionCode]; found {
			// Update if description changed
			if existing.Description != permission.Description {
				existing.Description = permission.Description
				if err = tx.Save(&existing).Error; err != nil {
					return err
				}
				log.Printf("updated permission: %s", permissionCode)
			}
		} else {
			// Create new permission
			if err = tx.Create(&permission).Error; err != nil {
				return err
			}
			log.Printf("created permission: %s", permissionCode)
		}
	}

	// Commit transaction
	return tx.Commit().Error
}

// SyncPermissions synchronizes permissions with the database
// It creates new permissions, updates existing ones, and optionally deactivates removed ones
func (s *PermissionSeeder) SyncPermissions() (err error) {
	log.Print("synchronizing permissions")

	// Begin transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Defer rollback in case of error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("panic occurred during permission sync: %v", r)
			}
		} else if err != nil {
			tx.Rollback()
		}
	}()

	// Get existing permissions
	var existingPermissions []entity.Permission
	if err = tx.Find(&existingPermissions).Error; err != nil {
		return err
	}

	// Create map for quick lookup
	existingMap := make(map[string]entity.Permission)
	for _, perm := range existingPermissions {
		existingMap[perm.Code] = perm
	}

	// Get all defined permission codes
	definedCodes := make(map[string]bool)
	for _, code := range permissions.GetAllPermissionCodes() {
		definedCodes[code] = true
	}

	// Process permissions
	for _, permission := range permissions.ApiPermissions {
		permissionCode := fmt.Sprintf("%s/%s", permission.Object, permission.Action)
		permission.Code = permissionCode

		// Check if permission exists
		if existing, found := existingMap[permissionCode]; found {
			// Update if description changed
			if existing.Description != permission.Description ||
				existing.IsDomainSpecific != permission.IsDomainSpecific {
				existing.Description = permission.Description
				existing.IsDomainSpecific = permission.IsDomainSpecific
				if err = tx.Save(&existing).Error; err != nil {
					return err
				}
				log.Printf("updated permission: %s", permissionCode)
			}
		} else {
			// Create new permission
			if err = tx.Create(&permission).Error; err != nil {
				return err
			}
			log.Printf("created permission: %s", permissionCode)
		}
	}

	// Commit transaction
	return tx.Commit().Error
}
