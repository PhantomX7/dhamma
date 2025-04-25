package middleware

import (
	"time"

	"github.com/PhantomX7/dhamma/constants"
	"github.com/PhantomX7/dhamma/utility/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a middleware function that logs request details, handles request IDs,
// sets up a context-specific logger, and logs any errors encountered during the request lifecycle.
// It also handle the error response if error occurred.
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

		// --- Error Logging ---
		// Check for errors after handlers have run
		requestErrors := c.Errors // Get errors attached to the context by handlers
		errorFields := []zap.Field{}
		if len(requestErrors) > 0 {
			// Log each error attached to the context
			for _, err := range requestErrors {
				// Optionally collect error messages for the final log entry
				errorFields = append(errorFields, zap.String("error", err.Error()))
			}
		}

		// --- End Error Logging ---

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
		if len(errorFields) > 0 {
			finalLogFields = append(finalLogFields, zap.Array("errors", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
				for _, err := range requestErrors {
					enc.AppendString(err.Error())
				}
				return nil
			})))
		}

		// Use appropriate log level based on status/errors
		if c.Writer.Status() >= 500 {
			contextLogger.Error("request completed with errors", finalLogFields...)
		} else if c.Writer.Status() >= 400 {
			contextLogger.Warn("request completed with client error", finalLogFields...)
		} else {
			contextLogger.Info("request completed successfully", finalLogFields...)
		}
	}
}
