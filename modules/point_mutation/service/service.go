package service

import (
	"github.com/PhantomX7/dhamma/modules/point_mutation"
)

type service struct {
	pointMutationRepo point_mutation.Repository
}

func New(
	pointMutationRepo point_mutation.Repository,
) point_mutation.Service {
	return &service{
		pointMutationRepo: pointMutationRepo,
	}
}
