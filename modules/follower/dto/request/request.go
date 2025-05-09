package request

import "github.com/PhantomX7/dhamma/utility/pagination"

type FollowerCreateRequest struct {
	DomainID uint64  `json:"domain_id" form:"domain_id" binding:"required,exist=domains.id"`
	Name     string  `json:"name" form:"name" binding:"required"`
	Phone    *string `json:"phone" form:"phone" binding:"omitempty"`
	IsYouth  *bool   `json:"is_muda_mudi" form:"is_muda_mudi" binding:"omitempty"`
}

type FollowerUpdateRequest struct {
	Name    *string `json:"name" form:"name" binding:"omitempty"`
	Phone   *string `json:"phone" form:"phone" binding:"omitempty"`
	IsYouth *bool   `json:"is_muda_mudi" form:"is_muda_mudi" binding:"omitempty"`
}

func NewFollowerPagination(conditions map[string][]string) *pagination.Pagination {
	filterDef := pagination.NewFilterDefinition().
		AddFilter("domain_id", pagination.FilterConfig{
			Field: "domain_id",
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
		AddFilter("phone", pagination.FilterConfig{
			Field: "phone",
			Type:  pagination.FilterTypeString,
			Operators: []pagination.FilterOperator{
				pagination.OperatorIn, pagination.OperatorEquals, pagination.OperatorLike,
			},
		}).
		AddFilter("is_youth", pagination.FilterConfig{
			Field: "is_youth",
			Type:  pagination.FilterTypeBool,
			Operators: []pagination.FilterOperator{
				pagination.OperatorEquals,
			},
		}).
		AddSort("name", pagination.SortConfig{
			Field:   "name",
			Allowed: true,
		}).
		AddSort("phone", pagination.SortConfig{
			Field:   "phone",
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
