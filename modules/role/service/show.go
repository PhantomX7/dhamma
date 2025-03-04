package service

import (
	"context"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/PhantomX7/dhamma/entity"
)

// Show implements role.Service
func (s *service) Show(roleID uint64, ctx context.Context) (role entity.Role, err error) {
	hasDomain, domainID := utility.GetDomainIDFromContext(ctx)

	role, err = s.roleRepo.FindByID(roleID, ctx)
	if err != nil {
		return
	}

	if hasDomain {
		if role.DomainID != domainID {
			return entity.Role{}, utility.LogError("you are not allowed to create role for another domain", nil)
		}
	}

	return
}
