package request

import "github.com/PhantomX7/dhamma/utility/pagination"

type FollowerCreateRequest struct {
	DomainID uint64  `json:"domain_id" form:"domain_id" binding:"required,exist=domains.id"`
	Name     string  `json:"name" form:"name" binding:"required"`
	Phone    *string `json:"phone" form:"phone" binding:"omitempty"`
	IsYouth  *bool   `json:"is_youth" form:"is_youth" binding:"omitempty"`
}

type FollowerUpdateRequest struct {
	Name    *string `json:"name" form:"name" binding:"omitempty"`
	Phone   *string `json:"phone" form:"phone" binding:"omitempty"`
	IsYouth *bool   `json:"is_youth" form:"is_youth" binding:"omitempty"`
}

// FollowerAddCardRequest defines the payload for adding a card to a follower.
type FollowerAddCardRequest struct {
	Code string `json:"code" form:"code" binding:"required,max=100,unique=cards.code"`
}

func NewFollowerPagination(conditions map[string][]string) *pagination.Pagination {
	filterDef := pagination.NewFilterDefinition().
		AddFilter("search", pagination.FilterConfig{
			TableName:    "followers",
			SearchFields: []string{"name", "phone"},
			Type:         pagination.FilterTypeString,
			Operators: []pagination.FilterOperator{
				pagination.OperatorLike,
			},
		}).
		// Add filter for domain name
		AddFilter("domain_name", pagination.FilterConfig{
			TableName: "Domain",
			Field:     "name",
			Type:      pagination.FilterTypeString,
			Operators: []pagination.FilterOperator{
				pagination.OperatorLike, pagination.OperatorEquals, pagination.OperatorIn,
			},
		}).
		AddFilter("domain_id", pagination.FilterConfig{
			TableName: "Domain",
			Field:     "id",
			Type:      pagination.FilterTypeID,
			Operators: []pagination.FilterOperator{
				pagination.OperatorEquals,
			},
		}).
		AddFilter("card_code", pagination.FilterConfig{
			TableName: "Card",
			Field:     "code",
			Type:      pagination.FilterTypeString,
			Operators: []pagination.FilterOperator{
				pagination.OperatorLike, pagination.OperatorEquals, pagination.OperatorIn,
			},
		}).
		AddFilter("name", pagination.FilterConfig{
			TableName: "followers",
			Field:     "name",
			Type:      pagination.FilterTypeString,
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
		AddFilter("created_at", pagination.FilterConfig{
			Field: "created_at",
			Type:  pagination.FilterTypeDate,
			Operators: []pagination.FilterOperator{
				pagination.OperatorBetween, pagination.OperatorEquals,
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
			DefaultOrder: "followers.id desc",
		},
	)
}
