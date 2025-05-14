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
	// Find the role by ID
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return
	}

	// Perform domain context check using DomainID from the request and the generic helper
	_, err = utility.CheckDomainContext(ctx, role.DomainID, "role", "delete permissions")
	if err != nil {
		return
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
