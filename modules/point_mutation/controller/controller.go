package controller

import (
	"github.com/PhantomX7/dhamma/modules/point_mutation"
)

type controller struct {
	pointMutationService point_mutation.Service
}

func New(pointMutationService point_mutation.Service) point_mutation.Controller {
	return &controller{
		pointMutationService: pointMutationService,
	}
}
