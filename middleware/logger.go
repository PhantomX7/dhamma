package middleware

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/PhantomX7/dhamma/constants"
	"github.com/PhantomX7/dhamma/utility/errors"
	"github.com/PhantomX7/dhamma/utility/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// LoggerConfig holds configuration for the logger middleware
type LoggerConfig struct {
	SkipPaths            []string      // Paths to skip logging (e.g., health checks)
	LogRequestBody       bool          // Whether to log request body
	LogResponseBody      bool          // Whether to log response body
	MaxBodySize          int64         // Maximum body size to log (bytes)
	SensitiveHeaders     []string      // Headers to redact in logs
	SlowRequestThreshold time.Duration // Threshold for slow request warnings
}

// DefaultLoggerConfig returns default configuration for logger middleware
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		SkipPaths:            []string{"/health", "/metrics", "/ping"},
		LogRequestBody:       false,
		LogResponseBody:      false,
		MaxBodySize:          1024 * 10, // 10KB
		SensitiveHeaders:     []string{"authorization", "cookie", "x-api-key", "x-auth-token"},
		SlowRequestThreshold: 2 * time.Second,
	}
}

// Logger is a middleware function that logs request details with enhanced structured logging,
// request tracing, performance metrics, and configurable options.
func (m *Middleware) Logger() gin.HandlerFunc {
	return m.LoggerWithConfig(DefaultLoggerConfig())
}

// LoggerWithConfig creates a logger middleware with custom configuration
func (m *Middleware) LoggerWithConfig(config LoggerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip logging for specified paths
		for _, skipPath := range config.SkipPaths {
			if c.Request.URL.Path == skipPath {
				c.Next()
				return
			}
		}

		start := time.Now()

		// Generate or get request ID and trace ID
		requestID := getOrGenerateRequestID(c)
		traceID := getOrGenerateTraceID(c)

		// Set IDs in context and headers
		c.Set(constants.EnumRequestIDKey, requestID)
		c.Set("trace_id", traceID)
		c.Header("X-Request-ID", requestID)
		c.Header("X-Trace-ID", traceID)

		// Extract user context if available
		userID := extractUserID(c)
		userRole := extractUserRole(c)

		// Create enhanced context logger
		contextLogger := createContextLogger(c, requestID, traceID, userID, userRole)

		// Add logger to context
		c.Request = c.Request.WithContext(logger.WithCtx(c.Request.Context(), contextLogger))

		// Log request start with optional body
		logRequestStart(contextLogger, c, config)

		// Capture response body if configured
		var responseBody *bytes.Buffer
		if config.LogResponseBody {
			responseBody = captureResponseBody(c)
		}

		// Process request (executes downstream handlers)
		c.Next()

		// Calculate duration and performance metrics
		duration := time.Since(start)
		performanceMetrics := calculatePerformanceMetrics(duration, c.Writer.Size())

		// Log request completion with enhanced details
		logRequestCompletion(contextLogger, c, duration, performanceMetrics, responseBody, config)
	}
}

// getOrGenerateRequestID gets request ID from header or generates a new one
func getOrGenerateRequestID(c *gin.Context) string {
	requestID := c.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}
	return requestID
}

// getOrGenerateTraceID gets trace ID from header or generates a new one
func getOrGenerateTraceID(c *gin.Context) string {
	traceID := c.GetHeader("X-Trace-ID")
	if traceID == "" {
		traceID = uuid.New().String()
	}
	return traceID
}

// extractUserID extracts user ID from context if available
func extractUserID(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			return id
		}
		if id, ok := userID.(int); ok {
			return strconv.Itoa(id)
		}
	}
	return ""
}

// extractUserRole extracts user role from context if available
func extractUserRole(c *gin.Context) string {
	if userRole, exists := c.Get("user_role"); exists {
		if role, ok := userRole.(string); ok {
			return role
		}
	}
	return ""
}

// createContextLogger creates an enhanced logger with request context
func createContextLogger(c *gin.Context, requestID, traceID, userID, userRole string) *zap.Logger {
	fields := []zap.Field{
		zap.String("request_id", requestID),
		zap.String("trace_id", traceID),
		zap.String("client_ip", c.ClientIP()),
		zap.String("real_ip", getRealIP(c)),
		zap.String("forwarded_for", c.GetHeader("X-Forwarded-For")),
	}

	if userID != "" {
		fields = append(fields, zap.String("user_id", userID))
	}
	if userRole != "" {
		fields = append(fields, zap.String("user_role", userRole))
	}

	return logger.Get().With(fields...)
}

// getRealIP gets the real IP address from various headers
func getRealIP(c *gin.Context) string {
	// Check X-Real-IP header first
	if realIP := c.GetHeader("X-Real-IP"); realIP != "" {
		return realIP
	}

	// Check X-Forwarded-For header
	if forwardedFor := c.GetHeader("X-Forwarded-For"); forwardedFor != "" {
		// Take the first IP from the comma-separated list
		if ips := strings.Split(forwardedFor, ","); len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Fallback to ClientIP
	return c.ClientIP()
}

// logRequestStart logs the beginning of request processing
func logRequestStart(contextLogger *zap.Logger, c *gin.Context, config LoggerConfig) {
	fields := []zap.Field{
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("query", c.Request.URL.RawQuery),
		zap.String("user_agent", c.Request.UserAgent()),
		zap.String("referer", c.Request.Referer()),
		zap.String("protocol", c.Request.Proto),
		zap.String("host", c.Request.Host),
		zap.Int64("content_length", c.Request.ContentLength),
		zap.String("content_type", c.Request.Header.Get("Content-Type")),
	}

	// Add sanitized headers
	headers := sanitizeHeaders(c.Request.Header, config.SensitiveHeaders)
	fields = append(fields, zap.Any("headers", headers))

	// Add request body if configured and within size limit
	if config.LogRequestBody && c.Request.ContentLength > 0 && c.Request.ContentLength <= config.MaxBodySize {
		if body := readRequestBody(c); body != "" {
			fields = append(fields, zap.String("request_body", body))
		}
	}

	contextLogger.Info("request started", fields...)
}

// readRequestBody safely reads and restores request body
func readRequestBody(c *gin.Context) string {
	if c.Request.Body == nil {
		return ""
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return ""
	}

	// Restore the body for downstream handlers
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	return string(body)
}

// captureResponseBody captures response body for logging
func captureResponseBody(c *gin.Context) *bytes.Buffer {
	responseBody := &bytes.Buffer{}
	writer := &responseBodyWriter{
		ResponseWriter: c.Writer,
		body:           responseBody,
	}
	c.Writer = writer
	return responseBody
}

// responseBodyWriter captures response body while writing
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseBodyWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

// calculatePerformanceMetrics calculates various performance metrics
func calculatePerformanceMetrics(duration time.Duration, responseSize int) map[string]interface{} {
	return map[string]interface{}{
		"duration_ms":      duration.Milliseconds(),
		"duration_seconds": duration.Seconds(),
		"response_size":    responseSize,
		"throughput_bps":   calculateThroughput(responseSize, duration),
		"is_slow_request":  duration > 2*time.Second,
	}
}

// calculateThroughput calculates bytes per second
func calculateThroughput(size int, duration time.Duration) float64 {
	if duration.Seconds() == 0 {
		return 0
	}
	return float64(size) / duration.Seconds()
}

// logRequestCompletion logs the completion of request processing
func logRequestCompletion(contextLogger *zap.Logger, c *gin.Context, duration time.Duration,
	performanceMetrics map[string]interface{}, responseBody *bytes.Buffer, config LoggerConfig) {

	status := c.Writer.Status()

	fields := []zap.Field{
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.Int("status", status),
		zap.String("status_text", getStatusText(status)),
		zap.Duration("duration", duration),
		zap.Any("performance", performanceMetrics),
		zap.Int("response_size", c.Writer.Size()),
	}

	// Add response body if configured and within size limit
	if config.LogResponseBody && responseBody != nil && int64(responseBody.Len()) <= config.MaxBodySize {
		fields = append(fields, zap.String("response_body", responseBody.String()))
	}

	// Add error information if present
	if len(c.Errors) > 0 {
		errorInfo := extractErrorInfo(c.Errors)
		fields = append(fields, zap.Any("errors", errorInfo))
	}

	// Log with appropriate level based on status and performance
	message := "request completed"
	if duration > config.SlowRequestThreshold {
		message = "slow request completed"
	}

	switch {
	case status >= 500:
		contextLogger.Error(message, fields...)
	case status >= 400:
		contextLogger.Warn(message, fields...)
	case duration > config.SlowRequestThreshold:
		contextLogger.Warn(message, fields...)
	default:
		contextLogger.Info(message, fields...)
	}
}

// sanitizeHeaders removes sensitive headers from logging
func sanitizeHeaders(headers map[string][]string, sensitiveHeaders []string) map[string]interface{} {
	sanitized := make(map[string]interface{})

	for key, values := range headers {
		lowerKey := strings.ToLower(key)
		isSensitive := false

		for _, sensitive := range sensitiveHeaders {
			if strings.ToLower(sensitive) == lowerKey {
				isSensitive = true
				break
			}
		}

		if isSensitive {
			sanitized[key] = "[REDACTED]"
		} else {
			if len(values) == 1 {
				sanitized[key] = values[0]
			} else {
				sanitized[key] = values
			}
		}
	}

	return sanitized
}

// extractErrorInfo extracts structured error information
func extractErrorInfo(ginErrors []*gin.Error) []map[string]interface{} {
	errorInfo := make([]map[string]interface{}, 0, len(ginErrors))

	for _, ginErr := range ginErrors {
		errorData := map[string]interface{}{
			"type": ginErr.Type,
			"meta": ginErr.Meta,
		}

		switch e := ginErr.Err.(type) {
		case *errors.AppError:
			errorData["error_type"] = e.Type
			errorData["error_code"] = e.Code
			errorData["message"] = e.Message
			errorData["status"] = e.Status
			if e.Err != nil {
				errorData["underlying_error"] = e.Err.Error()
			}
		case validator.ValidationErrors:
			errorData["error_type"] = "validation_error"
			errorData["validation_errors"] = formatValidationErrors(e)
		default:
			errorData["error_type"] = "unknown_error"
			errorData["message"] = e.Error()
		}

		errorInfo = append(errorInfo, errorData)
	}

	return errorInfo
}

// formatValidationErrors formats validation errors for logging
func formatValidationErrors(validationErrors validator.ValidationErrors) []map[string]string {
	formatted := make([]map[string]string, 0, len(validationErrors))

	for _, err := range validationErrors {
		formatted = append(formatted, map[string]string{
			"field": err.Field(),
			"tag":   err.Tag(),
			"value": err.Value().(string),
		})
	}

	return formatted
}

// getStatusText returns human-readable status text
func getStatusText(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "success"
	case status >= 300 && status < 400:
		return "redirect"
	case status >= 400 && status < 500:
		return "client_error"
	case status >= 500:
		return "server_error"
	default:
		return "unknown"
	}
}
