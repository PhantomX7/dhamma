package domain

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/domain/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Repository interface {
	Create(domain *entity.Domain, tx *gorm.DB, ctx context.Context) error
	Update(domain *entity.Domain, tx *gorm.DB, ctx context.Context) error
	FindAll(pg *pagination.Pagination, ctx context.Context) ([]entity.Domain, error)
	FindByID(domainID uint64, ctx context.Context) (entity.Domain, error)
	FindByCode(code string, ctx context.Context) (entity.Domain, error)
	GetDomainRoles(domainID uint64, ctx context.Context) (entity.Domain, error)
	Count(pg *pagination.Pagination, ctx context.Context) (int64, error)
}

type Service interface {
	Index(pg *pagination.Pagination, ctx context.Context) ([]entity.Domain, utility.PaginationMeta, error)
	Show(domainID uint64, ctx context.Context) (entity.Domain, error)
	Update(domainID uint64, request request.DomainUpdateRequest, ctx context.Context) (entity.Domain, error)
	Create(request request.DomainCreateRequest, ctx context.Context) (entity.Domain, error)
	ShowWithRoles(domainID uint64, ctx context.Context) (entity.Domain, error)
}

type Controller interface {
	Index(ctx *gin.Context)
	Show(ctx *gin.Context)
	Update(ctx *gin.Context)
	Create(ctx *gin.Context)
	//ShowWithRoles(ctx *gin.Context)
}
