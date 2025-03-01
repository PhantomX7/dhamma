package service

import (
	"context"
	"errors"
	"github.com/PhantomX7/dhamma/utility"
	"strings"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/user/dto/request"
)

func (u *service) Create(request request.UserCreateRequest, ctx context.Context) (user entity.User, err error) {
	haveDomain, domainID := utility.GetDomainIDFromContext(ctx)

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

	tx := u.transactionManager.NewTransaction()

	err = u.userRepo.Create(&user, tx, ctx)
	if err != nil {
		tx.Rollback()
		return
	}

	if haveDomain {
		err = u.userDomainRepo.AssignDomain(user.ID, domainID, tx, ctx)
		if err != nil {
			tx.Rollback()
			return
		}
	}

	tx.Commit()

	return
}
