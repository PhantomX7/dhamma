package service

import (
	"context"
	"net/http"

	"github.com/PhantomX7/dhamma/modules/role/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
)

// DeletePermissions removes specified permissions from a role.
func (s *service) DeletePermissions(ctx context.Context, roleID uint64, request request.RoleDeletePermissionsRequest) (err error) {
	// Get value from context
	contextValues, err := utility.ValuesFromContext(ctx)
	if err != nil {
		return err
	}

	// Find the role by ID
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return
	}

	// Validate domain access if context has a domain ID
	if contextValues.DomainID != nil {
		// Check if domain id in request is same as domain id in context
		if role.DomainID != *contextValues.DomainID {
			err = &errors.AppError{
				Message: "you are not allowed to delete permissions to role for another domain",
				Status:  http.StatusBadRequest,
			}
			return
		}
	}

	// Call casbin library to remove the permissions
	// Note: This assumes a DeletePermissions method exists in the casbin client.
	// We will need to implement this method next.
	err = s.casbin.DeleteRolePermissions(role.ID, role.DomainID, request.Permissions)
	if err != nil {
		err = &errors.AppError{
			Message: "fail to delete permissions",
			Status:  http.StatusBadRequest,
			Err:     err,
		}
		return
	}

	return nil // Return nil on success
}
