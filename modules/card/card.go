package card

import (
	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility/repository"
)

type Repository interface {
	repository.BaseRepositoryInterface[entity.Card]
}
