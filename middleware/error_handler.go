package middleware

import (
	"net/http"

	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
	"github.com/PhantomX7/dhamma/utility/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
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
			switch e := err.(type) {
			case *errors.AppError:
				// Handle application error
				c.AbortWithStatusJSON(e.Status, utility.BuildResponseFailed(e.Message, nil))

				return
			case validator.ValidationErrors:
				// Handle validation errors
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, utility.ValidationErrorResponse(e))

				return
			}

			logger.Get().Error("error happened", zap.Error(err))

			// Handle generic errors
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				utility.BuildResponseFailed("Internal server error", "An unexpected error occurred"),
			)
		}
	}
}
