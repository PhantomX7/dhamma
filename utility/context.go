package utility

import (
	"context"
	"errors"
	"github.com/PhantomX7/dhamma/constants"
	"github.com/gin-gonic/gin"
)

type ContextValues struct {
	DomainID *uint64
	UserID   uint64
	IsRoot   bool
}

// GetIDFromContext this is not a heavyweight operation
// it accesses payload from map in gin context
func GetIDFromContext(c *gin.Context) uint64 {
	id, _ := c.Get(constants.EnumJwtKeyUserId)
	return id.(uint64)
}

func GetDomainIDFromContext(context context.Context) (haveDomain bool, domainID uint64) {
	domainID, _ = context.Value(constants.EnumContextKeyDomainID).(uint64)

	if domainID == 0 {
		return
	}

	haveDomain = true
	return
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
