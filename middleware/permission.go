package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/PhantomX7/dhamma/constants/permissions"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/logger"
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

		// --- Permission Check Logic ---
		permissionCode := fmt.Sprintf("%s/%s", config.Object, config.Action)
		_, exists := m.permissionDefinitions[permissionCode]

		if !exists {
			// This should ideally not happen if permissions are seeded correctly
			// Log this occurrence for investigation
			logger.Get().Error("Permission definition not found in cache", zap.String("code", permissionCode))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "permission definition error",
			})
			return
		}

		var hasPermission bool
		var casbinErr error

		// Domain ID must be present for domain-specific permission check
		// if contextValues.DomainID == nil {
		// 	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		// 		"error": "domain context required for this permission",
		// 	})
		// 	return
		// }

		// Check domain-specific permission
		hasPermission, casbinErr = m.casbin.GetEnforcer().Enforce(
			fmt.Sprintf("%d", contextValues.UserID),
			fmt.Sprintf("%d", *contextValues.DomainID),
			config.Object,
			config.Action,
			permissions.PermissionTypeApi,
		)

		// Handle potential Casbin errors
		if casbinErr != nil {
			logger.Get().Error("Casbin enforce error", zap.Error(casbinErr))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "permission check error",
			})
			return
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "insufficient permissions",
			})
			return
		}

		c.Next()
	}
}
