package service

import (
	"context"
	"net/http"

	"github.com/PhantomX7/dhamma/modules/user/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
)

func (s *service) RemoveRole(ctx context.Context, userID uint64, request request.RemoveRoleRequest) (err error) {
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
			return &errors.AppError{
				Status:  http.StatusBadRequest,
				Message: "you are not allowed to remove role for another domain",
			}
		}
	}

	// check if user has this role domain
	hasDomain, err := s.userDomainRepo.HasDomain(ctx, userID, role.DomainID)
	if err != nil {
		return
	}
	if !hasDomain {
		return &errors.AppError{
			Status:  http.StatusBadRequest,
			Message: "you are not allowed to remove role for another domain",
		}
	}

	// Check if user already has the role
	hasRole, err := s.userRoleRepo.HasRole(ctx, userID, request.RoleID)
	if err != nil {
		return err
	}
	if !hasRole {
		return &errors.AppError{
			Status:  http.StatusBadRequest,
			Message: "user don't have the role",
		}
	}

	// Assign role in database
	err = s.userRoleRepo.RemoveRole(ctx, userID, role.ID, nil)
	if err != nil {
		return
	}

	// Assign role in casbin
	err = s.casbin.RemoveUserRole(user.ID, role.ID, role.DomainID)
	if err != nil {
		return
	}

	return
}
