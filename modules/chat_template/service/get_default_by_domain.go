package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
)

// GetDefaultByDomain implements chat_template.Service
func (s *service) GetDefaultByDomain(ctx context.Context, domainID uint64) (template entity.ChatTemplate, err error) {
	_, err = utility.CheckDomainContext(ctx, domainID, "chat template", "get default")
	if err != nil {
		return
	}

	template, err = s.chatTemplateRepo.GetDefaultByDomain(ctx, domainID)
	if err != nil {
		return
	}

	return
}
