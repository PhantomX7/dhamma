package utility

import (
	"github.com/PhantomX7/dhamma/constants"
	"github.com/gin-gonic/gin"
)

// GetRoleFromContext this is not a heavyweight operation
// it accesses payload from map in gin context
func GetRoleFromContext(c *gin.Context) string {
	role, _ := c.Get(constants.EnumJwtKeyRole)
	return role.(string)
}

// GetIDFromContext this is not a heavyweight operation
// it accesses payload from map in gin context
func GetIDFromContext(c *gin.Context) uint64 {
	id, _ := c.Get(constants.EnumJwtKeyUserId)
	return id.(uint64)
}
