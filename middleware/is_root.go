package middleware

import (
	"github.com/PhantomX7/dhamma/utility"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *Middleware) IsRoot() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role from context
		contextValues, err := utility.ValuesFromContext(c.Request.Context())
		if err != nil {
			return
		}

		// check if user is root
		if !contextValues.IsRoot {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "insufficient permissions",
			})
			return
		}

		c.Next()
	}
}
