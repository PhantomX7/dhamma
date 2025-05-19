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
		AddFilter("domain_name", pagination.FilterConfig{
			TableName: "Domain",
			Field:     "name",
			Type:      pagination.FilterTypeString,
			Operators: []pagination.FilterOperator{
				pagination.OperatorLike, pagination.OperatorEquals, pagination.OperatorIn,
			},
		}).
		AddFilter("code", pagination.FilterConfig{
			Field: "code",
			Type:  pagination.FilterTypeString,
			Operators: []pagination.FilterOperator{
				pagination.OperatorIn, pagination.OperatorEquals, pagination.OperatorLike,
			},
		}).
		AddFilter("description", pagination.FilterConfig{
			Field: "description",
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
			Field:     "created_at",
			Type:      pagination.FilterTypeDateTime,
			Operators: []pagination.FilterOperator{pagination.OperatorBetween},
		}).
		AddSort("name", pagination.SortConfig{
			Field:   "name",
			Allowed: true,
		}).
		AddSort("code", pagination.SortConfig{
			Field:   "code",
			Allowed: true,
		}).
		AddSort("description", pagination.SortConfig{
			Field:   "description",
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
