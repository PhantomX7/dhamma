package service

import (
	"context"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/chat_template/dto/request"
	"github.com/PhantomX7/dhamma/utility"
)

func (s *service) Update(ctx context.Context, templateID uint64, req request.ChatTemplateUpdateRequest) (template entity.ChatTemplate, err error) {
	template, err = s.chatTemplateRepo.FindByID(ctx, templateID)
	if err != nil {
		return
	}

	_, err = utility.CheckDomainContext(ctx, template.DomainID, "chat template", "update")
	if err != nil {
		return
	}

	// Store original default status
	originalIsDefault := template.IsDefault

	err = copier.Copy(&template, &req)
	if err != nil {
		return
	}

	// If IsDefault is being changed to true, handle the logic
	if req.IsDefault != nil && *req.IsDefault && !originalIsDefault {
		err = s.transactionManager.ExecuteInTransaction(func(tx *gorm.DB) error {
			// Update the template first
			if err = s.chatTemplateRepo.Update(ctx, &template, tx); err != nil {
				return err
			}

			// Set it as default (this will unset others)
			return s.chatTemplateRepo.SetAsDefault(ctx, template.ID, template.DomainID, tx)
		})
	} else {
		err = s.chatTemplateRepo.Update(ctx, &template, nil)
	}

	if err != nil {
		return
	}

	return
}
