package service

import (
	"context"
	"errors"

	"github.com/PhantomX7/dhamma/modules/user/dto/request"
	"github.com/PhantomX7/dhamma/utility"
)

func (s *service) AssignRole(ctx context.Context, userID uint64, request request.AssignRoleRequest) (err error) {
	// Get value from context
	contextValues, err := utility.ValuesFromContext(ctx)
	if err != nil {
		return
	}

	// Check if user exists
	user, err := s.userRepo.FindByID(ctx, userID, "Domains", "UserRoles.Role")
	if err != nil {
		return
	}

	role, err := s.roleRepo.FindByID(ctx, request.RoleID)
	if err != nil {
		return
	}

	// Check domain access if domain ID is set in context
	if contextValues.DomainID != nil {
		if role.DomainID != *contextValues.DomainID {
			return errors.New("you are not allowed to assign role for another domain")
		}
	}

	// Check if user already has the role
	hasRole, err := s.userRoleRepo.HasRole(ctx, userID, request.RoleID)
	if err != nil {
		return err
	}
	if hasRole {
		return errors.New("user already has this role")
	}

	// Assign role in database
	err = s.userRoleRepo.AssignRole(ctx, userID, role.DomainID, role.ID, nil)
	if err != nil {
		return
	}

	// Assign role in casbin
	err = s.casbin.AddUserRole(user.ID, role.ID, role.DomainID)
	if err != nil {
		return
	}

	return
}
