package service

import (
	"github.com/PhantomX7/dhamma/libs/transaction_manager"
	"github.com/PhantomX7/dhamma/modules/chat_template"
)

type service struct {
	chatTemplateRepo   chat_template.Repository
	transactionManager transaction_manager.Client
}

// New creates a new chat template service instance.
func New(
	chatTemplateRepo chat_template.Repository,
	transactionManager transaction_manager.Client,
) chat_template.Service {
	return &service{
		chatTemplateRepo:   chatTemplateRepo,
		transactionManager: transactionManager,
	}
}
