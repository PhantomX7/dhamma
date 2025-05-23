package request

import "github.com/PhantomX7/dhamma/utility/pagination"

type EventAttendanceCreateRequest struct {
	Name        string `json:"name" form:"name" binding:"required,unique=event_attendances.name"`
	Code        string `json:"code" form:"code" binding:"required,unique=event_attendances.code"`
	Description string `json:"description" form:"description"`
	IsActive    bool   `json:"is_active" form:"is_active" binding:"required"`
}

type EventAttendanceUpdateRequest struct {
	Name        *string `json:"name" form:"name" binding:"omitempty,unique=event_attendances.name"`
	Code        *string `json:"code" form:"code" binding:"omitempty,unique=event_attendances.code"`
	Description *string `json:"description" form:"description" binding:"omitempty"`
	IsActive    *bool   `json:"is_active" form:"is_active" binding:"omitempty"`
}

func NewEventAttendancePagination(conditions map[string][]string) *pagination.Pagination {
	filterDef := pagination.NewFilterDefinition().
		AddFilter("event_id", pagination.FilterConfig{
			Field: "event_id",
			Type:  pagination.FilterTypeID,
			Operators: []pagination.FilterOperator{
				pagination.OperatorIn, pagination.OperatorEquals,
			},
		}).
		AddFilter("follower_id", pagination.FilterConfig{
			Field: "follower_id",
			Type:  pagination.FilterTypeID,
			Operators: []pagination.FilterOperator{
				pagination.OperatorIn, pagination.OperatorEquals,
			},
		}).
		AddFilter("follower_name", pagination.FilterConfig{
			TableName: "Follower",
			Field:     "name",
			Type:      pagination.FilterTypeString,
			Operators: []pagination.FilterOperator{
				pagination.OperatorLike, pagination.OperatorEquals, pagination.OperatorIn,
			},
		}).
		AddFilter("created_at", pagination.FilterConfig{
			Field:     "created_at",
			Type:      pagination.FilterTypeDateTime,
			Operators: []pagination.FilterOperator{pagination.OperatorBetween, pagination.OperatorEquals},
		}).
		AddSort("created_at", pagination.SortConfig{
			Field:   "created_at",
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
