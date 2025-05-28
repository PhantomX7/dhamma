package request

import "github.com/PhantomX7/dhamma/utility/pagination"

type ChatTemplateCreateRequest struct {
	DomainID    uint64  `json:"domain_id" form:"domain_id" binding:"required,exist=domains.id"`
	Name        string  `json:"name" form:"name" binding:"required"`
	Description *string `json:"description" form:"description"`
	Content     string  `json:"content" form:"content" binding:"required"`
	IsDefault   *bool   `json:"is_default" form:"is_default"`
}

type ChatTemplateUpdateRequest struct {
	Name        *string `json:"name" form:"name" binding:"omitempty"`
	Description *string `json:"description" form:"description" binding:"omitempty"`
	Content     *string `json:"content" form:"content" binding:"omitempty"`
	IsDefault   *bool   `json:"is_default" form:"is_default"`
}

func NewChatTemplatePagination(conditions map[string][]string) *pagination.Pagination {
	filterDef := pagination.NewFilterDefinition().
		AddFilter("name", pagination.FilterConfig{
			TableName: "chat_templates",
			Field:     "name",
			Type:      pagination.FilterTypeString,
			Operators: []pagination.FilterOperator{
				pagination.OperatorIn, pagination.OperatorEquals, pagination.OperatorLike,
			},
		}).
		AddFilter("is_default", pagination.FilterConfig{
			TableName: "chat_templates",
			Field:     "is_default",
			Type:      pagination.FilterTypeBool,
			Operators: []pagination.FilterOperator{pagination.OperatorEquals},
		}).
		AddFilter("created_at", pagination.FilterConfig{
			TableName: "chat_templates",
			Field:     "created_at",
			Type:      pagination.FilterTypeDateTime,
			Operators: []pagination.FilterOperator{pagination.OperatorBetween, pagination.OperatorEquals},
		}).
		AddSort("name", pagination.SortConfig{
			TableName: "chat_templates",
			Field:     "name",
			Allowed:   true,
		}).
		AddSort("created_at", pagination.SortConfig{
			Field:   "created_at",
			Allowed: true,
		})

	return pagination.NewPagination(
		conditions,
		filterDef,
		pagination.PaginationOptions{
			DefaultLimit: 10,
			MaxLimit:     100,
		},
	)
}
