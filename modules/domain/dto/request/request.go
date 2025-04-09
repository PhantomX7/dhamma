package request

import "github.com/PhantomX7/dhamma/utility/pagination"

type DomainCreateRequest struct {
	Name        string `json:"name" form:"name" binding:"required,unique=domains.name"`
	Code        string `json:"code" form:"code" binding:"required,unique=domains.code"`
	Description string `json:"description" form:"description"`
	IsActive    *bool  `json:"is_active" form:"is_active" binding:"required"`
}

type DomainUpdateRequest struct {
	Name        *string `json:"name" form:"name" binding:"omitempty,unique=domains.name"`
	Code        *string `json:"code" form:"code" binding:"omitempty,unique=domains.code"`
	Description *string `json:"description" form:"description"`
	IsActive    *bool   `json:"is_active" form:"is_active"`
}

func NewDomainPagination(conditions map[string][]string) *pagination.Pagination {
	filterDef := pagination.NewFilterDefinition().
		AddFilter("id", pagination.FilterConfig{
			Field: "id",
			Type:  pagination.FilterTypeID,
			Operators: []pagination.FilterOperator{
				pagination.OperatorIn, pagination.OperatorEquals,
			},
		}).
		AddFilter("name", pagination.FilterConfig{
			Field: "name",
			Type:  pagination.FilterTypeString,
			Operators: []pagination.FilterOperator{
				pagination.OperatorIn, pagination.OperatorEquals, pagination.OperatorLike,
			},
		}).
		AddFilter("code", pagination.FilterConfig{
			Field: "code",
			Type:  pagination.FilterTypeString,
			Operators: []pagination.FilterOperator{
				pagination.OperatorIn, pagination.OperatorEquals, pagination.OperatorLike,
			},
		}).
		AddFilter("created_at", pagination.FilterConfig{
			Field:     "created_at",
			Type:      pagination.FilterTypeDateTime,
			Operators: []pagination.FilterOperator{pagination.OperatorBetween},
		}).
		AddFilter("is_active", pagination.FilterConfig{
			Field: "is_active",
			Type:  pagination.FilterTypeBool,
			Operators: []pagination.FilterOperator{
				pagination.OperatorEquals,
			},
		}).
		AddSort("id", pagination.SortConfig{
			Field:   "id",
			Allowed: true,
		}).
		AddSort("code", pagination.SortConfig{
			Field:   "code",
			Allowed: true,
		}).
		AddSort("name", pagination.SortConfig{
			Field:   "name",
			Allowed: true,
		})

	return pagination.NewPagination(
		conditions,
		filterDef,
		pagination.PaginationOptions{
			DefaultLimit: 20,
			MaxLimit:     100,
			DefaultOrder: "id desc",
		},
	)

}
