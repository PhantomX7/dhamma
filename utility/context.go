package utility

import (
	"context"
	"github.com/PhantomX7/dhamma/constants"
	"github.com/gin-gonic/gin"
)

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
