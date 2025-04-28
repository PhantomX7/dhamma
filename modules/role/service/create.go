package service

import (
	"context"
	"net/http"

	"github.com/jinzhu/copier"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/role/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
)

func (s *service) Create(ctx context.Context, request request.RoleCreateRequest) (role entity.Role, err error) {
	// Get value from context
	contextValues, err := utility.ValuesFromContext(ctx)
	if err != nil {
		return
	}

	// Check if domain id is set in context
	if contextValues.DomainID != nil {
		// Check if domain id in request is same as domain id in context
		if request.DomainID != *contextValues.DomainID {
			err = &errors.AppError{
				Message: "you cannot create role for other domain",
				Status:  http.StatusBadRequest,
			}
			return
		}
	}

	role = entity.Role{
		IsActive: true,
	}

	exist, err := s.roleRepo.Exists(ctx, map[string]any{
		"name":      request.Name,
		"domain_id": request.DomainID,
	})
	if err != nil {
		return
	}

	if exist {
		err = &errors.AppError{
			Message: "Role already exists",
			Status:  http.StatusBadRequest,
		}
		return
	}

	err = copier.Copy(&role, &request)
	if err != nil {
		return
	}

	err = s.roleRepo.Create(ctx, &role, nil)
	if err != nil {
		return
	}

	return
}
