package service

import (
	"context"
	"errors"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/PhantomX7/dhamma/entity"
)

// Show implements role.Service
func (s *service) Show(ctx context.Context, roleID uint64) (role entity.Role, err error) {
	// Get value from context
	contextValues, err := utility.ValuesFromContext(ctx)
	if err != nil {
		return
	}

	role, err = s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return
	}

	if contextValues.DomainID != nil {
		if role.DomainID != *contextValues.DomainID {
			return entity.Role{}, errors.New("you are not allowed to create role for another domain")
		}
	}

	return
}
