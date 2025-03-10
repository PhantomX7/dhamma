package middleware

import (
	"fmt"
	"net/http"

	"github.com/PhantomX7/dhamma/constants"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/gin-gonic/gin"
)

type PermissionConfig struct {
	Object string
	Action string
}

func (m *Middleware) Permission(config PermissionConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get value from context
		contextValues, err := utility.ValuesFromContext(c.Request.Context())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "context values not found",
			})
			return
		}

		// Root user bypass permission check
		if contextValues.IsRoot {
			c.Next()
			return
		}

		// Domain ID must be present for permission check
		if contextValues.DomainID == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "domain not found",
			})
			return
		}

		// Check if user has permission
		hasPermission, _ := m.casbin.GetEnforcer().Enforce(
			fmt.Sprintf("%d", contextValues.UserID),
			fmt.Sprintf("%d", *contextValues.DomainID),
			config.Object,
			config.Action,
			constants.EnumPermissionTypeApi,
		)

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "insufficient permissions",
			})
			return
		}

		c.Next()
	}
}
