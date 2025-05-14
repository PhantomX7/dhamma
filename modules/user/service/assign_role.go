package service

import (
	"context"
	"net/http"

	"github.com/PhantomX7/dhamma/modules/user/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
)

func (s *service) AssignRole(ctx context.Context, userID uint64, request request.AssignRoleRequest) (err error) {
	// Check if user exists
	user, err := s.userRepo.FindByID(ctx, userID, "Domains", "UserRoles.Role")
	if err != nil {
		return
	}

	role, err := s.roleRepo.FindByID(ctx, request.RoleID)
	if err != nil {
		return
	}

	// Perform domain context check using DomainID from the request and the generic helper
	_, err = utility.CheckDomainContext(ctx, role.DomainID, "user", "assign role")
	if err != nil {
		return
	}

	// check if user has this role domain
	hasDomain, err := s.userDomainRepo.HasDomain(ctx, userID, role.DomainID)
	if err != nil {
		return
	}
	if !hasDomain {
		return &errors.AppError{
			Status:  http.StatusBadRequest,
			Message: "you are not allowed to assign role for another domain",
		}
	}

	// Check if user already has the role
	hasRole, err := s.userRoleRepo.HasRole(ctx, userID, request.RoleID)
	if err != nil {
		return err
	}
	if hasRole {
		return &errors.AppError{
			Status:  http.StatusBadRequest,
			Message: "user already has the role",
		}
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
