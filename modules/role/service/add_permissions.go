package service

import (
	"context"
	"net/http"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/role/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
)

func (s *service) AddPermissions(ctx context.Context, roleID uint64, request request.RoleAddPermissionsRequest) (role entity.Role, err error) {
	// Get value from context
	contextValues, err := utility.ValuesFromContext(ctx)
	if err != nil {
		return
	}

	role, err = s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return
	}

	// Check if domain id is set in context
	if contextValues.DomainID != nil {
		// Check if domain id in request is same as domain id in context
		if role.DomainID != *contextValues.DomainID {
			err = &errors.AppError{
				Message: "you are not allowed to add permissions to role for another domain",
				Status:  http.StatusBadRequest,
			}
			return
		}
	}

	s.casbin.AddRolePermissions(role.ID, role.DomainID, request.Permissions)

	return
}
