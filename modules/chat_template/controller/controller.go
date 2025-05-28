package controller

import (
	"github.com/PhantomX7/dhamma/modules/chat_template"
)

type controller struct {
	chatTemplateService chat_template.Service
}

func New(chatTemplateService chat_template.Service) chat_template.Controller {
	return &controller{
		chatTemplateService: chatTemplateService,
	}
}
