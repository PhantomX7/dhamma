package request

import "github.com/PhantomX7/dhamma/utility/pagination"

type RoleCreateRequest struct {
	DomainID    uint64 `json:"domain_id" form:"domain_id" binding:"required,exist=domains.id"`
	Name        string `json:"name" form:"name" binding:"required,unique=roles.name"`
	Description string `json:"description" form:"description"`
	IsActive    *bool  `json:"is_active" form:"is_active" binding:"required"`
}

type RoleUpdateRequest struct {
	Name        *string `json:"name" form:"name" binding:"unique=roles.name"`
	Description *string `json:"description" form:"description"`
	IsActive    *bool   `json:"is_active" form:"is_active"`
}

func NewRolePagination(conditions map[string][]string) *pagination.Pagination {
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
		AddFilter("created_at", pagination.FilterConfig{
			Field:     "created_at",
			Type:      pagination.FilterTypeDateTime,
			Operators: []pagination.FilterOperator{pagination.OperatorBetween},
		}).
		AddSort("id", pagination.SortConfig{
			Field:   "id",
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
