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

func (s *service) Create(request request.RoleCreateRequest, ctx context.Context) (role entity.Role, err error) {
	hasDomain, domainID := utility.GetDomainIDFromContext(ctx)

	if hasDomain {
		if request.DomainID != domainID {
			return entity.Role{}, utility.LogError("you are not allowed to create role for another domain", nil)
		}
	}

	role = entity.Role{
		IsActive: true,
	}

	existingRole, err := s.roleRepo.FindByNameAndDomainID(request.Name, request.DomainID, ctx)
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

	err = s.roleRepo.Create(&role, nil, ctx)
	if err != nil {
		return
	}

	return
}
