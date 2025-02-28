package controller

import (
	"github.com/PhantomX7/dhamma/modules/domain"
)

type controller struct {
	domainService domain.Service
}

func New(domainService domain.Service) domain.Controller {
	return &controller{
		domainService: domainService,
	}
}
