package chat_template

import (
	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/chat_template/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
	"github.com/PhantomX7/dhamma/utility/repository"
)

type Repository interface {
	repository.BaseRepositoryInterface[entity.ChatTemplate]
	SetAsDefault(ctx context.Context, templateID uint64, domainID uint64, tx *gorm.DB) error
	GetDefaultByDomain(ctx context.Context, domainID uint64) (entity.ChatTemplate, error)
}

type Service interface {
	Index(ctx context.Context, paginationConfig *pagination.Pagination) ([]entity.ChatTemplate, utility.PaginationMeta, error)
	Show(ctx context.Context, templateID uint64) (entity.ChatTemplate, error)
	Create(ctx context.Context, request request.ChatTemplateCreateRequest) (entity.ChatTemplate, error)
	Update(ctx context.Context, templateID uint64, request request.ChatTemplateUpdateRequest) (entity.ChatTemplate, error)
	SetAsDefault(ctx context.Context, templateID uint64) (entity.ChatTemplate, error)
	GetDefaultByDomain(ctx context.Context, domainID uint64) (entity.ChatTemplate, error)
}

type Controller interface {
	Index(c *gin.Context)
	Show(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	SetAsDefault(c *gin.Context)
	GetDefaultByDomain(c *gin.Context)
}
