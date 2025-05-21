package middleware

import (
	"time"

	"github.com/PhantomX7/dhamma/constants"
	"github.com/PhantomX7/dhamma/utility/errors"
	"github.com/PhantomX7/dhamma/utility/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Logger is a middleware function that logs request details, handles request IDs,
// sets up a context-specific logger, and logs any errors encountered during the request lifecycle.
func (m *Middleware) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate or get request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set request ID in context and header
		c.Set(constants.EnumRequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)

		// Add logger with request ID to context
		contextLogger := logger.Get().With(
			zap.String("request_id", requestID),
			zap.String("client_ip", c.ClientIP()),
		)

		// Add logger to context
		c.Request = c.Request.WithContext(logger.WithCtx(c.Request.Context(), contextLogger))
		// Process request (executes downstream handlers)
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Update Prometheus metrics
		// metrics.HttpRequestsTotal.WithLabelValues(
		// 	c.Request.Method,
		// 	c.Request.URL.Path,
		// 	//string(rune(c.Writer.Status())), // Status might not be final yet, consider moving if needed
		// 	//requestID,
		// ).Inc()

		// metrics.HttpRequestDuration.WithLabelValues(
		// 	c.Request.Method,
		// 	c.Request.URL.Path,
		// 	requestID,
		// ).Observe(duration.Seconds())

		// Log request completion
		// Add error fields to the final log if any occurred
		finalLogFields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.Int("body_size", c.Writer.Size()),
		}

		// Use appropriate log level based on status/errors
		if c.Writer.Status() >= 500 {
			contextLogger.Error("request completed with errors", finalLogFields...)
		} else if c.Writer.Status() >= 400 {
			// Check if there are any errors
			if len(c.Errors) > 0 {
				// Get the last error
				err := c.Errors.Last().Err

				// Handle different error types
				switch e := err.(type) {
				case *errors.AppError:
					// Handle application error
					finalLogFields = append(finalLogFields, zap.Error(e.Err))
					contextLogger.Warn("request completed with client error", finalLogFields...)
				case validator.ValidationErrors:
					// Handle validation errors
					finalLogFields = append(finalLogFields, zap.String("error", "validation_errors"))
				default:
					// Handle other types of errors
					finalLogFields = append(finalLogFields, zap.Error(e))
					contextLogger.Warn("request completed with client error", finalLogFields...)
				}
			} else {
				contextLogger.Warn("request completed with client error", finalLogFields...)
			}
		} else {
			contextLogger.Info("request completed successfully", finalLogFields...)
		}
	}
}
