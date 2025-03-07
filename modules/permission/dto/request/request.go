package request

import "github.com/PhantomX7/dhamma/utility/pagination"

type PermissionCreateRequest struct {
	Name        string `json:"name" form:"name" binding:"required"`
	Object      string `json:"object" form:"object" binding:"required"`
	Action      string `json:"action" form:"action" binding:"required"`
	Description string `json:"description" form:"description"`
}

type PermissionUpdateRequest struct {
	Name        *string `json:"name" form:"name"`
	Description *string `json:"description" form:"description"`
}

func NewPermissionPagination(conditions map[string][]string) *pagination.Pagination {
	filterDef := pagination.NewFilterDefinition().
		AddFilter("id", pagination.FilterConfig{
			Field: "id",
			Type:  pagination.FilterTypeID,
			Operators: []pagination.FilterOperator{
				pagination.OperatorIn, pagination.OperatorEquals,
			},
		}).
		AddFilter("code", pagination.FilterConfig{
			Field: "code",
			Type:  pagination.FilterTypeString,
			Operators: []pagination.FilterOperator{
				pagination.OperatorIn, pagination.OperatorEquals, pagination.OperatorLike,
			},
		}).
		AddSort("id", pagination.SortConfig{
			Field:   "id",
			Allowed: true,
		}).
		AddSort("code", pagination.SortConfig{
			Field:   "code",
			Allowed: true,
		})

	return pagination.NewPagination(
		conditions,
		filterDef,
		pagination.PaginationOptions{
			DefaultLimit: 1000,
			MaxLimit:     1000,
			DefaultOrder: "code asc",
		},
	)

}
