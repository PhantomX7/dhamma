package request

import "github.com/PhantomX7/dhamma/utility/pagination"

type UserCreateRequest struct {
	Username string `form:"username" json:"username" binding:"required,unique=users.username"`
	Password string `form:"password" json:"password" binding:"required"`
}
type AssignDomainRequest struct {
	DomainID uint64 `json:"domain_id" form:"domain_id" binding:"required,exist=domains.id"`
}

type AssignRoleRequest struct {
	RoleID uint64 `json:"role_id" form:"role_id" binding:"required,exist=roles.id"`
}

func NewUserPagination(conditions map[string][]string) *pagination.Pagination {
	filterDef := pagination.NewFilterDefinition().
		AddFilter("id", pagination.FilterConfig{
			Field: "id",
			Type:  pagination.FilterTypeID,
			Operators: []pagination.FilterOperator{
				pagination.OperatorIn, pagination.OperatorEquals,
			},
		}).
		AddFilter("username", pagination.FilterConfig{
			Field: "username",
			Type:  pagination.FilterTypeString,
			Operators: []pagination.FilterOperator{
				pagination.OperatorIn, pagination.OperatorEquals, pagination.OperatorLike,
			},
		}).
		AddFilter("is_active", pagination.FilterConfig{
			Field: "is_active",
			Type:  pagination.FilterTypeBool,
			Operators: []pagination.FilterOperator{
				pagination.OperatorEquals,
			},
		}).
		AddFilter("created_at", pagination.FilterConfig{
			Field: "created_at",
			Type:  pagination.FilterTypeDate,
			Operators: []pagination.FilterOperator{
				pagination.OperatorBetween, pagination.OperatorEquals,
			},
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
