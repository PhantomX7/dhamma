package service

import (
	"github.com/PhantomX7/dhamma/modules/domain"
)

type service struct {
	domainRepo domain.Repository
}

func New(
	domainRepo domain.Repository,
) domain.Service {
	return &service{
		domainRepo: domainRepo,
	}
}
