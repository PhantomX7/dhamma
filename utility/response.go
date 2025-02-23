package utility

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/stoewer/go-strcase"
)

type PaginationMeta struct {
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Total  int64 `json:"total"`
}

type Response struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Error   any            `json:"error,omitempty"`
	Data    any            `json:"data,omitempty"`
	Meta    PaginationMeta `json:"meta,omitempty"`
}

func BuildPaginationResponseSuccess(message string, data any, meta PaginationMeta) Response {
	res := Response{
		Status:  true,
		Message: message,
		Error:   nil,
		Data:    data,
	}
	return res
}

func BuildResponseSuccess(message string, data any) Response {
	res := Response{
		Status:  true,
		Message: message,
		Error:   nil,
		Data:    data,
	}
	return res
}

func BuildResponseFailed(message string, err any) Response {
	res := Response{
		Status:  false,
		Message: message,
		Error:   err,
		Data:    nil,
	}
	return res
}

func ValidationErrorResponse(err error) Response {
	errs, ok := err.(validator.ValidationErrors)
	if ok {
		list := make(map[string]string)
		for _, err := range errs {
			list[strcase.SnakeCase(err.Field())] = formatValidationError(err)
		}
		return BuildResponseFailed("Validation error", list)
	}

	return BuildResponseFailed("Validation error", err.Error())
}

func formatValidationError(e validator.FieldError) string {
	field := strcase.SnakeCase(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "exists":
		return fmt.Sprintf("%s is required", field)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, e.Param())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, e.Param())
	case "email":
		return "Invalid email format"
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, e.Param())
	case "unique":
		return fmt.Sprintf("%s must be unique", field)
	case "exist":
		return fmt.Sprintf("%s does not exist", field)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, e.Param())
	case "date":
		return fmt.Sprintf("%s must be in date format", field)
	case "value":
		return fmt.Sprintf(
			"%s must be one of %s",
			field,
			strings.Join(strings.Split(e.Param(), "."), ", "),
		)
	default:
		return fmt.Sprintf("%s is not valid", field)
	}
}
