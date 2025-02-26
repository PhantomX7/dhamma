package repository

import (
	"context"
	"errors"

	"github.com/PhantomX7/dhamma/entity"
)

// FindByUsername retrieves a user by their username.
//
// It takes a username as a string and a context for database operations.
// It returns a User and an error. If the error is not nil, it indicates that
// the query failed. If the error is nil, the user is successfully retrieved.
func (r *repository) FindByUsername(username string, ctx context.Context) (userM entity.User, err error) {

	err = r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&userM).Error
	if err != nil {
		err = errors.New("error find user by username")
		return
	}

	return
}
