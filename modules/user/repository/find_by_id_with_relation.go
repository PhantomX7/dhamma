package repository

import (
	"context"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/PhantomX7/dhamma/entity"
)

// FindByIDWithRelation finds a user by id and load all the relations.
//
// It finds a user by id and load Domains and UserRoles with Role and Domain
// relations.
//
// It returns a User and an error. If the error is not nil, it means that the
// query failed. If the error is nil, then the user is returned.
func (r *repository) FindByIDWithRelation(userID uint64, ctx context.Context) (userM entity.User, err error) {

	err = r.db.WithContext(ctx).
		Preload("Domains").
		Preload("UserRoles.Role").
		Preload("UserRoles.Domain").
		Where("id = ?", userID).Take(&userM).Error
	if err != nil {
		err = utility.LogError("error find user by id with relation", err)
		return
	}

	return
}
