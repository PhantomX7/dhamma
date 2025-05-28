package service

import (
	"context"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/chat_template/dto/request"
	"github.com/PhantomX7/dhamma/utility"
)

func (s *service) Create(ctx context.Context, req request.ChatTemplateCreateRequest) (template entity.ChatTemplate, err error) {
	_, err = utility.CheckDomainContext(ctx, req.DomainID, "chat_template", "create")
	if err != nil {
		return
	}

	err = copier.Copy(&template, &req)
	if err != nil {
		return
	}

	// Set default values
	template.IsActive = true
	if req.IsDefault != nil {
		template.IsDefault = *req.IsDefault
	} else {
		template.IsDefault = false
	}

	// If this template is being set as default, handle the logic
	if template.IsDefault {
		err = s.transactionManager.ExecuteInTransaction(func(tx *gorm.DB) error {
			// Create the template first
			if err = s.chatTemplateRepo.Create(ctx, &template, tx); err != nil {
				return err
			}

			// Set it as default (this will unset others)
			return s.chatTemplateRepo.SetAsDefault(ctx, template.ID, template.DomainID, tx)
		})
	} else {
		err = s.chatTemplateRepo.Create(ctx, &template, nil)
	}

	if err != nil {
		return
	}

	return
}
