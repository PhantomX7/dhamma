package middleware

import (
	"time"

	"github.com/PhantomX7/dhamma/constants"
	"github.com/PhantomX7/dhamma/metrics"
	"github.com/PhantomX7/dhamma/utility/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

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
		contextLogger := logger.Logger.With(
			zap.String("request_id", requestID),
			zap.String("client_ip", c.ClientIP()),
		)
		c.Set("logger", contextLogger)

		// Log request start
		contextLogger.Info("request started",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("user_agent", c.Request.UserAgent()),
		)

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Update Prometheus metrics
		metrics.HttpRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.Request.URL.Path,
			//string(rune(c.Writer.Status())),
			//requestID,
		).Inc()

		metrics.HttpRequestDuration.WithLabelValues(
			c.Request.Method,
			c.Request.URL.Path,
			requestID,
		).Observe(duration.Seconds())

		// Log request completion
		contextLogger.Info("request completed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.Int("body_size", c.Writer.Size()),
		)
	}
}
