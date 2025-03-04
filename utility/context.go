package utility

import (
	"context"
	"errors"
)

type ContextValues struct {
	DomainID *uint64
	UserID   uint64
	IsRoot   bool
}

func NewContextWithValues(ctx context.Context, values ContextValues) context.Context {
	return context.WithValue(ctx, "values", values)
}

func ValuesFromContext(ctx context.Context) (ContextValues, error) {
	values, ok := ctx.Value("values").(ContextValues)
	if !ok {
		return ContextValues{}, errors.New("context values not found")
	}
	return values, nil
}
