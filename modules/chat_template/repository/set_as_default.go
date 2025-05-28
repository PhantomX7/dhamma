package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility/errors"
	"gorm.io/gorm"
)

// SetAsDefault sets a chat template as default and unsets others in the same domain.
func (r *repository) SetAsDefault(ctx context.Context, templateID uint64, domainID uint64, tx *gorm.DB) error {
	// Start a transaction
	db := r.db
	if tx != nil {
		db = tx
	}

	// First, unset all default templates in the domain
	if err := db.Model(&entity.ChatTemplate{}).
		Where("domain_id = ?", domainID).
		Update("is_default", false).Error; err != nil {
		return errors.WrapError(err, "failed to unset default templates")
	}

	// Then, set the specified template as default
	if err := db.Model(&entity.ChatTemplate{}).
		Where("id = ? AND domain_id = ?", templateID, domainID).
		Update("is_default", true).Error; err != nil {
		return errors.WrapError(err, "failed to set template as default")
	}

	return nil
}
