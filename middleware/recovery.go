package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"

	"github.com/PhantomX7/dhamma/constants"
	"github.com/PhantomX7/dhamma/utility/errors"
	"github.com/PhantomX7/dhamma/utility/logger"
)

// RecoveryConfig holds configuration for the recovery middleware
type RecoveryConfig struct {
	EnableStackTrace     bool                                                        // Whether to include stack trace in logs
	EnableCircuitBreaker bool                                                        // Whether to enable circuit breaker
	MaxRequests          uint32                                                      // Maximum requests allowed to pass through when circuit breaker is half-open
	Interval             time.Duration                                               // Interval to clear the internal counts
	Timeout              time.Duration                                               // Timeout for circuit breaker
	ReadyToTrip          func(counts gobreaker.Counts) bool                          // Function to determine when to trip the circuit
	OnStateChange        func(name string, from gobreaker.State, to gobreaker.State) // Callback for state changes
}

// DefaultRecoveryConfig returns default configuration for recovery middleware
func DefaultRecoveryConfig() RecoveryConfig {
	return RecoveryConfig{
		EnableStackTrace:     true,
		EnableCircuitBreaker: true,
		MaxRequests:          10,
		Interval:             60 * time.Second,
		Timeout:              30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Trip if failure rate is over 60% and we have at least 5 requests
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 5 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			logger.Get().Warn("Circuit breaker state changed",
				zap.String("name", name),
				zap.String("from", from.String()),
				zap.String("to", to.String()),
			)
		},
	}
}

// CircuitBreakerManager manages multiple circuit breakers for different endpoints
type CircuitBreakerManager struct {
	breakers map[string]*gobreaker.CircuitBreaker
	mu       sync.RWMutex
	config   RecoveryConfig
}

// NewCircuitBreakerManager creates a new circuit breaker manager
func NewCircuitBreakerManager(config RecoveryConfig) *CircuitBreakerManager {
	return &CircuitBreakerManager{
		breakers: make(map[string]*gobreaker.CircuitBreaker),
		config:   config,
	}
}

// GetBreaker gets or creates a circuit breaker for the given name
func (cbm *CircuitBreakerManager) GetBreaker(name string) *gobreaker.CircuitBreaker {
	cbm.mu.RLock()
	breaker, exists := cbm.breakers[name]
	cbm.mu.RUnlock()

	if exists {
		return breaker
	}

	cbm.mu.Lock()
	defer cbm.mu.Unlock()

	// Double-check after acquiring write lock
	if breaker, exists = cbm.breakers[name]; exists {
		return breaker
	}

	// Create new circuit breaker
	settings := gobreaker.Settings{
		Name:          name,
		MaxRequests:   cbm.config.MaxRequests,
		Interval:      cbm.config.Interval,
		Timeout:       cbm.config.Timeout,
		ReadyToTrip:   cbm.config.ReadyToTrip,
		OnStateChange: cbm.config.OnStateChange,
	}

	breaker = gobreaker.NewCircuitBreaker(settings)
	cbm.breakers[name] = breaker

	return breaker
}

// GetStats returns statistics for all circuit breakers
func (cbm *CircuitBreakerManager) GetStats() map[string]gobreaker.Counts {
	cbm.mu.RLock()
	defer cbm.mu.RUnlock()

	stats := make(map[string]gobreaker.Counts)
	for name, breaker := range cbm.breakers {
		stats[name] = breaker.Counts()
	}

	return stats
}

// Global circuit breaker manager instance
var globalCBManager *CircuitBreakerManager
var cbManagerOnce sync.Once

// getCircuitBreakerManager returns the global circuit breaker manager instance
func getCircuitBreakerManager(config RecoveryConfig) *CircuitBreakerManager {
	cbManagerOnce.Do(func() {
		globalCBManager = NewCircuitBreakerManager(config)
	})
	return globalCBManager
}

// Recovery middleware handles panics and implements circuit breaker patterns
func (m *Middleware) Recovery() gin.HandlerFunc {
	return m.RecoveryWithConfig(DefaultRecoveryConfig())
}

// RecoveryWithConfig creates a recovery middleware with custom configuration
func (m *Middleware) RecoveryWithConfig(config RecoveryConfig) gin.HandlerFunc {
	cbManager := getCircuitBreakerManager(config)

	return func(c *gin.Context) {
		// Get request context information
		requestID := c.GetString(constants.EnumRequestIDKey)
		if requestID == "" {
			requestID = "unknown"
		}

		// Get context logger or fallback to global logger
		contextLogger := logger.FromCtx(c.Request.Context())
		if contextLogger == nil {
			contextLogger = logger.Get()
		}

		// Create endpoint identifier for circuit breaker
		endpointKey := fmt.Sprintf("%s:%s", c.Request.Method, c.Request.URL.Path)

		// Defer panic recovery
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with full context
				logPanic(contextLogger, c, err, config.EnableStackTrace)

				// Create error response
				appError := errors.NewInternalError(
					"Internal server error occurred",
					fmt.Errorf("panic recovered: %v", err),
				)

				// Set error in context for error handler middleware
				c.Error(appError)

				// If response hasn't been written, send error response
				if !c.Writer.Written() {
					c.JSON(http.StatusInternalServerError, gin.H{
						"success": false,
						"error": gin.H{
							"code":    appError.Code,
							"message": "Internal server error occurred",
							"type":    appError.Type,
						},
						"request_id": requestID,
					})
				}

				// Abort the request
				c.Abort()
			}
		}()

		// Execute request with circuit breaker if enabled
		if config.EnableCircuitBreaker {
			breaker := cbManager.GetBreaker(endpointKey)

			// Execute request through circuit breaker
			_, err := breaker.Execute(func() (interface{}, error) {
				c.Next()

				// Check if there were any errors or if status indicates failure
				if len(c.Errors) > 0 {
					return nil, fmt.Errorf("request failed with errors")
				}

				// Consider 5xx status codes as failures for circuit breaker
				if c.Writer.Status() >= 500 {
					return nil, fmt.Errorf("request failed with status %d", c.Writer.Status())
				}

				return nil, nil
			})

			// Handle circuit breaker errors
			if err != nil {
				if err == gobreaker.ErrOpenState {
					// Circuit breaker is open, return service unavailable
					contextLogger.Warn("Circuit breaker is open",
						zap.String("endpoint", endpointKey),
						zap.String("request_id", requestID),
					)

					appError := errors.NewServiceUnavailableError(
						"Service temporarily unavailable",
						nil,
					)

					c.Error(appError)

					if !c.Writer.Written() {
						c.JSON(http.StatusServiceUnavailable, gin.H{
							"success": false,
							"error": gin.H{
								"code":    appError.Code,
								"message": "Service temporarily unavailable due to high error rate",
								"type":    appError.Type,
							},
							"request_id":  requestID,
							"retry_after": int(config.Timeout.Seconds()),
						})
					}

					c.Abort()
					return
				} else if err == gobreaker.ErrTooManyRequests {
					// Too many requests in half-open state
					contextLogger.Warn("Circuit breaker: too many requests",
						zap.String("endpoint", endpointKey),
						zap.String("request_id", requestID),
					)

					appError := errors.NewTooManyRequestsError(
						"Too many requests",
						nil,
					)

					c.Error(appError)

					if !c.Writer.Written() {
						c.JSON(http.StatusTooManyRequests, gin.H{
							"success": false,
							"error": gin.H{
								"code":    appError.Code,
								"message": "Too many requests, please try again later",
								"type":    appError.Type,
							},
							"request_id": requestID,
						})
					}

					c.Abort()
					return
				}
				// Other errors are handled by the circuit breaker internally
			}
		} else {
			// Execute request without circuit breaker
			c.Next()
		}
	}
}

// logPanic logs panic information with full context
func logPanic(contextLogger *zap.Logger, c *gin.Context, panicErr interface{}, includeStackTrace bool) {
	fields := []zap.Field{
		zap.Any("panic", panicErr),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("query", c.Request.URL.RawQuery),
		zap.String("client_ip", c.ClientIP()),
		zap.String("user_agent", c.Request.UserAgent()),
		zap.String("referer", c.Request.Referer()),
	}

	// Add request ID if available
	if requestID := c.GetString(constants.EnumRequestIDKey); requestID != "" {
		fields = append(fields, zap.String("request_id", requestID))
	}

	// Add user context if available
	if userID := c.GetString("user_id"); userID != "" {
		fields = append(fields, zap.String("user_id", userID))
	}

	if domainID := c.GetString("domain_id"); domainID != "" {
		fields = append(fields, zap.String("domain_id", domainID))
	}

	// Add stack trace if enabled
	if includeStackTrace {
		fields = append(fields, zap.String("stack_trace", string(debug.Stack())))
	}

	contextLogger.Error("Panic recovered", fields...)
}

// GetCircuitBreakerStats returns statistics for all circuit breakers
func (m *Middleware) GetCircuitBreakerStats() map[string]gobreaker.Counts {
	if globalCBManager == nil {
		return make(map[string]gobreaker.Counts)
	}
	return globalCBManager.GetStats()
}

// ResetCircuitBreaker resets a specific circuit breaker
func (m *Middleware) ResetCircuitBreaker(endpointKey string) bool {
	if globalCBManager == nil {
		return false
	}

	globalCBManager.mu.RLock()
	_, exists := globalCBManager.breakers[endpointKey]
	globalCBManager.mu.RUnlock()

	if !exists {
		return false
	}

	// Create a new breaker to effectively reset it
	globalCBManager.mu.Lock()
	defer globalCBManager.mu.Unlock()

	settings := gobreaker.Settings{
		Name:          endpointKey,
		MaxRequests:   globalCBManager.config.MaxRequests,
		Interval:      globalCBManager.config.Interval,
		Timeout:       globalCBManager.config.Timeout,
		ReadyToTrip:   globalCBManager.config.ReadyToTrip,
		OnStateChange: globalCBManager.config.OnStateChange,
	}

	globalCBManager.breakers[endpointKey] = gobreaker.NewCircuitBreaker(settings)
	return true
}
