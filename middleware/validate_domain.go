package middleware

import (
	"net/http"

	"github.com/PhantomX7/dhamma/utility"
	"github.com/gin-gonic/gin"
)

func (m *Middleware) ValidateDomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get value from context
		contextValues, err := utility.ValuesFromContext(c.Request.Context())
		if err != nil {
			return
		}

		domainCode := c.Param("domain_code")
		domain, err := m.domainRepo.FindOneByField(c.Request.Context(), "code", domainCode)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "domain not found",
			})
			return
		}

		if !contextValues.IsRoot {
			// Check if user has domain
			hasDomain, err := m.userDomainRepo.HasDomain(c.Request.Context(), contextValues.UserID, domain.ID)
			if err != nil {
				return
			}

			if !hasDomain {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "insufficient permissions",
				})
				return
			}
		}

		c.Request = c.Request.WithContext(utility.NewContextWithValues(
			c.Request.Context(),
			utility.ContextValues{
				DomainID: &domain.ID,
				UserID:   contextValues.UserID,
				IsRoot:   contextValues.IsRoot,
			},
		))

		c.Next()
	}
}
