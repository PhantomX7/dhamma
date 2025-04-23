package service

import (
	"context"

	"github.com/PhantomX7/dhamma/modules/role/dto/request"
	"github.com/PhantomX7/dhamma/utility"
)

// DeletePermissions removes specified permissions from a role.
func (s *service) DeletePermissions(ctx context.Context, roleID uint64, request request.RoleDeletePermissionsRequest) (err error) {
	// Get value from context
	contextValues, err := utility.ValuesFromContext(ctx)
	if err != nil {
		return err // Return context error directly
	}

	// Find the role by ID
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		// Wrap error for consistent handling (e.g., not found)
		return utility.WrapError(err, "role not found")
	}

	// Validate domain access if context has a domain ID
	if contextValues.DomainID != nil {
		if role.DomainID != *contextValues.DomainID {
			// Return a specific error for permission denied
			return utility.WrapError(utility.ErrPermissionDenied, "you cannot delete other domain permission") // Or a more specific error
		}
	}

	// Call casbin library to remove the permissions
	// Note: This assumes a DeletePermissions method exists in the casbin client.
	// We will need to implement this method next.
	err = s.casbin.DeleteRolePermissions(role.ID, role.DomainID, request.Permissions)
	if err != nil {
		// Wrap potential casbin errors
		return utility.WrapError(err, "failed to delete permissions in casbin")
	}

	return nil // Return nil on success
}
