package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
)

// GetDefaultByDomain implements chat_template.Service
func (s *service) GetDefaultByDomain(ctx context.Context, domainID uint64) (template entity.ChatTemplate, err error) {
	template, err = s.chatTemplateRepo.GetDefaultByDomain(ctx, domainID)
	if err != nil {
		return
	}

	return
}
