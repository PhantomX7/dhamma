package service

import (
	"context"

	"github.com/PhantomX7/dhamma/modules/user/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"gorm.io/gorm"
)

func (s *service) RemoveDomain(ctx context.Context, userID uint64, request request.RemoveDomainRequest) error {
	// Check if user exists
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return utility.WrapError(utility.ErrNotFound, "user not found")
	}

	// Check if domain exists
	_, err = s.domainRepo.FindByID(ctx, request.DomainID)
	if err != nil {
		return utility.WrapError(utility.ErrNotFound, "domain not found")
	}

	// Check if user has the domain
	hasDomain, err := s.userDomainRepo.HasDomain(ctx, userID, request.DomainID)
	if err != nil {
		return err
	}

	if !hasDomain {
		return utility.WrapError(utility.ErrNotFound, "user is not assigned to this domain")
	}

	// Start transaction
	err = s.transactionManager.ExecuteInTransaction(func(tx *gorm.DB) error {
		// Remove domain from user
		err = s.userDomainRepo.RemoveDomain(ctx, userID, request.DomainID, tx)
		if err != nil {
			return err
		}

		// Remove all roles associated with this domain
		err = s.userRoleRepo.RemoveRolesByUserAndDomainID(ctx, userID, request.DomainID, tx)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
