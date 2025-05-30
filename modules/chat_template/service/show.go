package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
)

// Show implements chat_template.Service
func (s *service) Show(ctx context.Context, templateID uint64) (template entity.ChatTemplate, err error) {
	template, err = s.chatTemplateRepo.FindByID(ctx, templateID, "Domain")
	if err != nil {
		return
	}

	_, err = utility.CheckDomainContext(ctx, template.DomainID, "chat template", "show")
	if err != nil {
		return
	}

	return
}
