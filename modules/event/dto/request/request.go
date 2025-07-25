package request

import "github.com/PhantomX7/dhamma/utility/pagination"

type EventCreateRequest struct {
	DomainID      uint64 `json:"domain_id" form:"domain_id" binding:"required,exist=domains.id"`
	Name          string `json:"name" form:"name" binding:"required"`
	Description   string `json:"description" form:"description"`
	PointsAwarded int    `json:"points_awarded" form:"points_awarded" binding:"required"`
}

type EventUpdateRequest struct {
	Name          *string `json:"name" form:"name" binding:"omitempty"`
	Description   *string `json:"description" form:"description" binding:"omitempty"`
	PointsAwarded *int    `json:"points_awarded" form:"points_awarded" binding:"omitempty"`
}

// EventAttendRequest defines the payload for a follower attending an event.
type EventAttendRequest struct {
	CardCode string `json:"card_code" form:"card_code" binding:"required,exist=cards.code"`
}

// EventAttendRequest defines the payload for a follower attending an event.
type EventAttendByIDRequest struct {
	FollowerID uint64 `json:"follower_id" form:"follower_id" binding:"required,exist=followers.id"`
}

func NewEventPagination(conditions map[string][]string) *pagination.Pagination {
	filterDef := pagination.NewFilterDefinition().
		AddFilter("name", pagination.FilterConfig{
			TableName: "events",
			Field:     "name",
			Type:      pagination.FilterTypeString,
			Operators: []pagination.FilterOperator{
				pagination.OperatorIn, pagination.OperatorEquals, pagination.OperatorLike,
			},
		}).
		AddFilter("created_at", pagination.FilterConfig{
			TableName: "events",
			Field:     "created_at",
			Type:      pagination.FilterTypeDateTime,
			Operators: []pagination.FilterOperator{pagination.OperatorBetween, pagination.OperatorEquals},
		}).
		AddSort("name", pagination.SortConfig{
			TableName: "events",
			Field:     "name",
			Allowed:   true,
		}).
		AddSort("points_awarded", pagination.SortConfig{
			Field:   "points_awarded",
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

func NewEventAttendancePagination(conditions map[string][]string) *pagination.Pagination {
	filterDef := pagination.NewFilterDefinition().
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
		AddSort("name", pagination.SortConfig{
			Field:   "name",
			Allowed: true,
		}).
		AddSort("points_awarded", pagination.SortConfig{
			Field:   "points_awarded",
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
