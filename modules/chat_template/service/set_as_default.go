package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
)

// SetAsDefault implements chat_template.Service
func (s *service) SetAsDefault(ctx context.Context, templateID uint64) (template entity.ChatTemplate, err error) {
	template, err = s.chatTemplateRepo.FindByID(ctx, templateID)
	if err != nil {
		return
	}

	err = s.chatTemplateRepo.SetAsDefault(ctx, templateID, template.DomainID, nil)
	if err != nil {
		return
	}

	// Refresh the template to get updated default status
	template, err = s.chatTemplateRepo.FindByID(ctx, templateID)
	if err != nil {
		return
	}

	return
}
