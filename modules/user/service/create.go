package service

import (
	"context"
	"errors"
	"github.com/PhantomX7/dhamma/utility"
	"gorm.io/gorm"
	"strings"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/user/dto/request"
)

func (s *service) Create(ctx context.Context, request request.UserCreateRequest) (user entity.User, err error) {
	// Get value from context
	contextValues, err := utility.ValuesFromContext(ctx)
	if err != nil {
		return
	}

	//haveDomain, domainID := utility.GetDomainIDFromContext(ctx)

	user = entity.User{
		IsActive:     true,
		IsSuperAdmin: false,
	}
	request.Username = strings.ToLower(strings.TrimSpace(request.Username))

	_ = copier.Copy(&user, &request)

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		err = errors.New("failed to hash password")
		return
	}
	user.Password = string(password)

	err = s.transactionManager.ExecuteInTransaction(func(tx *gorm.DB) error {
		err = s.userRepo.Create(ctx, &user, tx)
		if err != nil {
			return err
		}

		if contextValues.DomainID != nil {
			err = s.userDomainRepo.AssignDomain(ctx, user.ID, *contextValues.DomainID, tx)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return
	}

	return
}
