package repository

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility/errors"
	"github.com/PhantomX7/dhamma/utility/logger"
)

func (r *repository) DeleteInvalidToken(ctx context.Context) error {
	err := r.db.WithContext(ctx).
		Where("expires_at < ? OR is_valid = ?", time.Now(), false).
		Delete(&entity.RefreshToken{}).Error

	if err != nil {
		errMessage := "failed to delete invalid refresh token"
		logger.FromCtx(ctx).Error(errMessage, zap.Error(err))
		return errors.WrapError(errors.ErrDatabase, errMessage)
	}

	return nil
}
