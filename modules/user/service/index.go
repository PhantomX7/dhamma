package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/go-core/utility/request_util"
)

// Index implements user.Service.
func (s *service) Index(config request_util.PaginationConfig, ctx context.Context) ([]entity.User, error) {
	panic("unimplemented")
}
