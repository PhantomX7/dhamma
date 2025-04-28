package service

import (
	"context"
	"net/http"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
)

// Show implements role.Service
func (s *service) Show(ctx context.Context, roleID uint64) (role entity.Role, err error) {
	// Get value from context
	contextValues, err := utility.ValuesFromContext(ctx)
	if err != nil {
		return
	}

	role, err = s.roleRepo.FindByID(ctx, roleID, "Domain")
	if err != nil {
		return
	}

	// Check if domain id is set in context
	if contextValues.DomainID != nil {
		// Check if domain id in request is same as domain id in context
		if role.DomainID != *contextValues.DomainID {
			err = &errors.AppError{
				Message: "you are not allowed to view role for another domain",
				Status:  http.StatusBadRequest,
			}
			return
		}
	}

	role.Permissions = s.casbin.GetRolePermissions(role.ID, role.DomainID)

	return
}
