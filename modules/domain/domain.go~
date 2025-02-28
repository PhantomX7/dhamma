package domain

import (
	"context"
	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(domain *entity.Domain, tx *gorm.DB, ctx context.Context) error
	Update(domain *entity.Domain, tx *gorm.DB, ctx context.Context) error
	FindAll(pg *pagination.Pagination, ctx context.Context) ([]entity.Domain, error)
	FindByID(domainID uint64, ctx context.Context) (entity.Domain, error)
	FindByName(name string, ctx context.Context) (entity.Domain, error)
	GetDomainRoles(domainID uint64, ctx context.Context) (entity.Domain, error)
	Count(pg *pagination.Pagination, ctx context.Context) (int64, error)
}
