package utility

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/stoewer/go-strcase"
)

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
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
			list[strcase.SnakeCase(err.Field())] = validationErrorToText(err)
		}
		return BuildResponseFailed("Validation error", list)
	}

	return BuildResponseFailed("Validation error", err.Error())

}

func validationErrorToText(e validator.FieldError) string {
	errorField := strcase.SnakeCase(e.Field())
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", errorField)
	case "exists":
		return fmt.Sprintf("%s is required", errorField)
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s character", errorField, e.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than %s character", errorField, e.Param())
	case "email":
		return "Invalid email format"
	case "len":
		return fmt.Sprintf("%s must be %s characters long", errorField, e.Param())
	case "unique":
		return fmt.Sprintf("%s already exist", errorField)
	case "exist":
		return fmt.Sprintf("%s does not exist", errorField)
	case "gte":
		return fmt.Sprintf("%s must equal to or greater than %s", errorField, e.Param())
	case "lte":
		return fmt.Sprintf("%s must equal to or less than %s", errorField, e.Param())
	case "date":
		return fmt.Sprintf("%s must be in date format", errorField)
	case "value":
		return fmt.Sprintf("%s value must be any of %s", errorField, strings.Join(strings.Split(e.Param(), "."), ", "))
	}
	return fmt.Sprintf("%s is not valid", errorField)
}
