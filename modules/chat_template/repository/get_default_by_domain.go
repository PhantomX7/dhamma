package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility/errors"
)

// GetDefaultByDomain gets the default chat template for a domain.
func (r *repository) GetDefaultByDomain(ctx context.Context, domainID uint64) (entity.ChatTemplate, error) {
	var template entity.ChatTemplate
	err := r.db.WithContext(ctx).
		Where("domain_id = ? AND is_default = ? AND is_active = ?", domainID, true, true).
		First(&template).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return template, errors.WrapError(err, "default chat template not found for domain")
		}
		return template, errors.WrapError(err, "failed to get default chat template")
	}

	return template, nil
}
