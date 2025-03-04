package service

import (
	"context"
	"errors"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/role/dto/request"
)

func (s *service) Create(ctx context.Context, request request.RoleCreateRequest) (role entity.Role, err error) {
	// Get value from context
	contextValues, err := utility.ValuesFromContext(ctx)
	if err != nil {
		return
	}

	if contextValues.DomainID != nil {
		if request.DomainID != *contextValues.DomainID {
			return entity.Role{}, errors.New("you are not allowed to create role for another domain")
		}
	}

	role = entity.Role{
		IsActive: true,
	}

	existingRole, err := s.roleRepo.FindByNameAndDomainID(ctx, request.Name, request.DomainID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	if existingRole.ID != 0 {
		err = errors.New("role with this name already exist for this domain")
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
