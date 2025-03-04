package repository

import (
	"context"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/PhantomX7/dhamma/entity"
)

// FindByUsername retrieves a user by their username.
//
// It takes a username as a string and a context for database operations.
// It returns a User and an error. If the error is not nil, it indicates that
// the query failed. If the error is nil, the user is successfully retrieved.
func (r *repository) FindByUsername(ctx context.Context, username string) (userM entity.User, err error) {

	err = r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&userM).Error
	if err != nil {
		err = utility.LogError("error find user by username", err)
		return
	}

	return
}
