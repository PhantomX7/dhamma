package middleware

import (
	"github.com/PhantomX7/dhamma/constants"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *Middleware) IsRoot() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role from context
		userRole, exists := c.Get(constants.EnumJwtKeyRole)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "you are not allowed to access this resource",
			})
			return
		}

		// Type assertion
		roleStr, ok := userRole.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "role is not in correct format",
			})
			return
		}

		// Check if role matches
		if roleStr != constants.EnumRoleRoot {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "insufficient permissions",
			})
			return
		}

		c.Next()
	}
}
