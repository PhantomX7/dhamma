package request

import "github.com/PhantomX7/dhamma/utility/pagination"

func NewUserPagination(conditions map[string][]string) *pagination.Pagination {
	filterDef := pagination.NewFilterDefinition().
		AddFilter("id", pagination.FilterConfig{
			Field: "id",
			Type:  pagination.FilterTypeID,
			Operators: []pagination.FilterOperator{
				pagination.OperatorIn, pagination.OperatorEquals,
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
		AddSort("username", pagination.SortConfig{
			Field:   "username",
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
