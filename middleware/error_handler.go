package middleware

import (
	"net/http"

	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
	"github.com/gin-gonic/gin"
)

// ErrorHandler middleware handles application errors
func (m *Middleware) ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			// Get the last error
			err := c.Errors.Last().Err

			// Check if response has already been written
			if c.Writer.Written() {
				return
			}

			// Handle different error types
			if appErr, ok := err.(*errors.AppError); ok {
				switch appErr.Type {
				case errors.ErrorTypeValidation:
					c.AbortWithStatusJSON(appErr.Status, utility.ValidationErrorResponse(appErr.Err))
				case errors.ErrorTypeService:
					c.AbortWithStatusJSON(appErr.Status, utility.BuildResponseFailed(appErr.Message, appErr.Err.Error()))
				case errors.ErrorTypeBadRequest:
					c.AbortWithStatusJSON(appErr.Status, utility.BuildResponseFailed(appErr.Message, appErr.Err.Error()))
				default:
					c.AbortWithStatusJSON(http.StatusInternalServerError, utility.BuildResponseFailed("Internal server error", "An unexpected error occurred"))
				}
				return
			}

			// Handle generic errors
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				utility.BuildResponseFailed("Internal server error", err.Error()),
			)
		}
	}
}
