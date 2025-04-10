package domain

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/domain/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
	"github.com/PhantomX7/dhamma/utility/repository"

	"github.com/gin-gonic/gin"
)

type Repository interface {
	repository.BaseRepositoryInterface[entity.Domain]
	FindByCode(ctx context.Context, code string) (entity.Domain, error)
	GetDomainRoles(ctx context.Context, domainID uint64) (entity.Domain, error)
}

type Service interface {
	Index(ctx context.Context, pg *pagination.Pagination) ([]entity.Domain, utility.PaginationMeta, error)
	Show(ctx context.Context, domainID uint64) (entity.Domain, error)
	Update(ctx context.Context, domainID uint64, request request.DomainUpdateRequest) (entity.Domain, error)
	Create(ctx context.Context, request request.DomainCreateRequest) (entity.Domain, error)
	ShowWithRoles(ctx context.Context, domainID uint64) (entity.Domain, error)
}

type Controller interface {
	Index(ctx *gin.Context)
	Show(ctx *gin.Context)
	Update(ctx *gin.Context)
	Create(ctx *gin.Context)
	//ShowWithRoles(ctx *gin.Context)
}
