package middleware

import (
	"context"
	"github.com/PhantomX7/dhamma/constants"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *Middleware) ValidateUserDomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role from context
		id, exists := c.Get(constants.EnumJwtKeyUserId)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "you are not allowed to access this resource",
			})
			return
		}

		// Type assertion
		userID, ok := id.(uint64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "role is not in correct format",
			})
			return
		}

		domainCode := c.Param("domain_code")
		domain, err := m.domainRepo.FindByCode(domainCode, c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "domain not found",
			})
			return
		}

		hasDomain, err := m.userDomainRepo.HasDomain(userID, domain.ID, c)
		if err != nil {
			return
		}

		if !hasDomain {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "insufficient permissions",
			})
			return
		}

		// Set user domain in context
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), constants.EnumContextKeyDomainID, domain.ID))

		c.Next()
	}
}
