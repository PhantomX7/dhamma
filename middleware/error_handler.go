package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/PhantomX7/dhamma/constants"
	"github.com/PhantomX7/dhamma/utility/errors"
	"github.com/PhantomX7/dhamma/utility/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stoewer/go-strcase"
	"go.uber.org/zap"
)

// ErrorHandler middleware handles application errors with enhanced logging and context
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

			// Get request context information
			requestID := c.GetString(constants.EnumRequestIDKey)
			userID := c.GetString("user_id")     // Assuming user_id is set in auth middleware
			domainID := c.GetString("domain_id") // Assuming domain_id is set in domain middleware

			// Get context logger or fallback to global logger
			contextLogger := logger.FromCtx(c.Request.Context())

			// Common log fields for all error types
			commonFields := []zap.Field{
				zap.String("request_id", requestID),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("client_ip", c.ClientIP()),
				zap.String("user_agent", c.Request.UserAgent()),
			}

			// Add user and domain context if available
			if userID != "" {
				commonFields = append(commonFields, zap.String("user_id", userID))
			}
			if domainID != "" {
				commonFields = append(commonFields, zap.String("domain_id", domainID))
			}

			// Handle different error types
			switch e := err.(type) {
			case *errors.AppError:
				// Enrich AppError with request context if not already set
				if e.RequestID == "" {
					e = e.WithRequestID(requestID)
				}

				// Log based on error type and severity
				errorFields := append(commonFields,
					zap.String("error_type", string(e.Type)),
					zap.String("error_code", e.Code),
					zap.String("error_message", e.Message),
				)

				// Add underlying error if present
				if e.Err != nil {
					errorFields = append(errorFields, zap.Error(e.Err))
				}

				// Add error details if present
				if e.Details != nil {
					errorFields = append(errorFields, zap.Any("error_details", e.Details))
				}

				// Log with appropriate level based on error type
				switch e.Type {
				case errors.ErrorTypeValidation, errors.ErrorTypeAuthentication, errors.ErrorTypeAuthorization:
					contextLogger.Warn("Client error occurred", errorFields...)
				case errors.ErrorTypeDatabase, errors.ErrorTypeExternal, errors.ErrorTypeInternal:
					contextLogger.Error("Server error occurred", errorFields...)
				case errors.ErrorTypeRateLimit:
					contextLogger.Warn("Rate limit exceeded", errorFields...)
				default:
					contextLogger.Info("Application error occurred", errorFields...)
				}

				// Build error response
				errorResponse := map[string]interface{}{
					"success": false,
					"error": map[string]interface{}{
						"message":    e.Message,
						"code":       e.Code,
						"type":       e.Type,
						"request_id": requestID,
					},
				}

				// Add details to response if present and not sensitive
				if e.Details != nil && !m.isSensitiveError(e.Type) {
					errorResponse["error"].(map[string]interface{})["details"] = e.Details
				}

				c.AbortWithStatusJSON(e.Status, errorResponse)
				return

			case validator.ValidationErrors:
				// Handle validation errors with detailed field information
				// validationFields := append(commonFields,
				// 	zap.String("error_type", "validation_error"),
				// 	zap.Int("validation_error_count", len(e)),
				// )

				// contextLogger.Warn("Validation error occurred", validationFields...)

				// Build detailed validation error response
				validationResponse := map[string]interface{}{
					"success": false,
					"error": map[string]interface{}{
						"message":    "Validation failed",
						"code":       errors.CodeValidationFailed,
						"type":       errors.ErrorTypeValidation,
						"request_id": requestID,
						"details":    m.formatValidationErrors(e),
					},
				}

				c.AbortWithStatusJSON(http.StatusBadRequest, validationResponse)
				return

			default:
				// Handle unexpected errors
				unexpectedFields := append(commonFields,
					zap.Error(err),
					zap.String("error_type", "unexpected_error"),
					zap.String("stack_trace", string(debug.Stack())),
				)

				contextLogger.Error("Unexpected error occurred", unexpectedFields...)

				// Create internal error response
				internalError := errors.NewInternalError("An unexpected error occurred", err).WithRequestID(requestID)

				errorResponse := map[string]interface{}{
					"success": false,
					"error": map[string]interface{}{
						"message":    "Internal server error",
						"code":       internalError.Code,
						"type":       internalError.Type,
						"request_id": requestID,
					},
				}

				c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse)
			}
		}
	}
}

// isSensitiveError determines if error details should be hidden from response
func (m *Middleware) isSensitiveError(errorType errors.ErrorType) bool {
	sensitiveTypes := []errors.ErrorType{
		errors.ErrorTypeDatabase,
		errors.ErrorTypeInternal,
		errors.ErrorTypeExternal,
	}

	for _, sensitiveType := range sensitiveTypes {
		if errorType == sensitiveType {
			return true
		}
	}
	return false
}

// formatValidationErrors converts validator.ValidationErrors to a structured format
func (m *Middleware) formatValidationErrors(validationErrors validator.ValidationErrors) map[string]interface{} {
	errorDetails := make(map[string]interface{})
	fieldErrors := make(map[string]string)

	for _, err := range validationErrors {

		fieldName := strcase.SnakeCase(err.Field())
		tag := err.Tag()
		_ = err.Value() // Acknowledge but don't use the value

		// Create human-readable error message based on validation tag
		var message string
		switch tag {
		case "required":
			message = "This field is required"
		case "email":
			message = "Must be a valid email address"
		case "min":
			message = "Value is too short or small"
		case "max":
			message = "Value is too long or large"
		case "len":
			message = "Value length is invalid"
		case "numeric":
			message = "Must be a number"
		case "alpha":
			message = "Must contain only letters"
		case "alphanum":
			message = "Must contain only letters and numbers"
		case "unique":
			message = "Value must be unique"
		case "exists":
			message = "Value does not exist"
		default:
			message = "Invalid value"
		}

		fieldErrors[fieldName] = message
	}

	errorDetails["fields"] = fieldErrors
	errorDetails["total_errors"] = len(validationErrors)

	return errorDetails
}
