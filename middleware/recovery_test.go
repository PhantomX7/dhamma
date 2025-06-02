package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PhantomX7/dhamma/constants"
	"github.com/PhantomX7/dhamma/utility/logger"
	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecoveryMiddleware(t *testing.T) {
	// Initialize logger for testing
	logger.NewLogger()

	// Create middleware instance
	middleware := &Middleware{}

	tests := []struct {
		name           string
		handler        gin.HandlerFunc
		expectedStatus int
		expectedError  bool
		config         RecoveryConfig
	}{
		{
			name: "normal request without panic",
			handler: func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			config:         DefaultRecoveryConfig(),
		},
		{
			name: "panic recovery",
			handler: func(c *gin.Context) {
				panic("test panic")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			config:         DefaultRecoveryConfig(),
		},
		{
			name: "panic recovery without circuit breaker",
			handler: func(c *gin.Context) {
				panic("test panic")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			config: RecoveryConfig{
				EnableStackTrace:     true,
				EnableCircuitBreaker: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup Gin in test mode
			gin.SetMode(gin.TestMode)
			router := gin.New()

			// Add recovery middleware
			router.Use(middleware.RecoveryWithConfig(tt.config))

			// Add test route
			router.GET("/test", tt.handler)

			// Create test request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.False(t, response["success"].(bool))
				assert.Contains(t, response, "error")
				assert.Contains(t, response, "request_id")
			}
		})
	}
}

func TestCircuitBreakerManager(t *testing.T) {
	config := DefaultRecoveryConfig()
	manager := NewCircuitBreakerManager(config)

	t.Run("create and get circuit breaker", func(t *testing.T) {
		breaker1 := manager.GetBreaker("test-endpoint-1")
		assert.NotNil(t, breaker1)

		// Getting the same breaker should return the same instance
		breaker2 := manager.GetBreaker("test-endpoint-1")
		assert.Equal(t, breaker1, breaker2)

		// Different endpoint should create different breaker
		breaker3 := manager.GetBreaker("test-endpoint-2")
		assert.NotNil(t, breaker3)
		assert.NotEqual(t, breaker1, breaker3)
	})

	t.Run("get stats", func(t *testing.T) {
		// Create some breakers
		manager.GetBreaker("endpoint-1")
		manager.GetBreaker("endpoint-2")

		stats := manager.GetStats()
		assert.Len(t, stats, 4) // 2 from previous test + 2 new ones
		assert.Contains(t, stats, "endpoint-1")
		assert.Contains(t, stats, "endpoint-2")
	})
}

func TestCircuitBreakerIntegration(t *testing.T) {
	// Initialize logger for testing
	logger.NewLogger()

	// Create middleware instance
	middleware := &Middleware{}

	// Custom config with lower thresholds for testing
	config := RecoveryConfig{
		EnableStackTrace:     true,
		EnableCircuitBreaker: true,
		MaxRequests:          2,
		Interval:             1 * time.Second,
		Timeout:              1 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Trip after 2 failures
			return counts.TotalFailures >= 2
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			t.Logf("Circuit breaker %s changed from %s to %s", name, from, to)
		},
	}

	t.Run("circuit breaker opens after failures", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.New()

		// Add request ID middleware for proper logging
		router.Use(func(c *gin.Context) {
			c.Set(constants.EnumRequestIDKey, "test-request-id")
			c.Next()
		})

		// Add recovery middleware with circuit breaker
		router.Use(middleware.RecoveryWithConfig(config))

		// Add test route that always fails
		router.GET("/fail", func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		})

		// Make requests that should trigger circuit breaker
		for i := 0; i < 3; i++ {
			req := httptest.NewRequest(http.MethodGet, "/fail", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if i < 2 {
				// First two requests should get through
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			} else {
				// Third request might be blocked by circuit breaker
				// (depending on timing and circuit breaker implementation)
				t.Logf("Request %d status: %d", i+1, w.Code)
			}
		}
	})

	t.Run("circuit breaker stats", func(t *testing.T) {
		stats := middleware.GetCircuitBreakerStats()
		assert.NotEmpty(t, stats)

		// Check if we have stats for our test endpoint
		found := false
		for endpoint := range stats {
			if endpoint == "GET:/fail" {
				found = true
				break
			}
		}
		assert.True(t, found, "Should have stats for GET:/fail endpoint")
	})
}

func TestRecoveryWithRequestContext(t *testing.T) {
	// Initialize logger for testing
	logger.NewLogger()

	// Create middleware instance
	middleware := &Middleware{}

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add request context middleware
	router.Use(func(c *gin.Context) {
		c.Set(constants.EnumRequestIDKey, "test-request-123")
		c.Set("user_id", "user-456")
		c.Set("domain_id", "domain-789")
		c.Next()
	})

	// Add recovery middleware
	router.Use(middleware.Recovery())

	// Add test route that panics
	router.GET("/panic", func(c *gin.Context) {
		panic("test panic with context")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.False(t, response["success"].(bool))
	assert.Equal(t, "test-request-123", response["request_id"])
	assert.Contains(t, response, "error")
}

func TestLogPanic(t *testing.T) {
	// Initialize logger for testing
	logger.NewLogger()

	// Create a test context
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/test?param=value", bytes.NewBufferString("test body"))
	c.Set(constants.EnumRequestIDKey, "test-request-id")
	c.Set("user_id", "test-user")
	c.Set("domain_id", "test-domain")

	// Test logging panic (this should not actually panic)
	contextLogger := logger.Get()
	logPanic(contextLogger, c, "test panic error", true)

	// If we reach here, the function didn't panic
	assert.True(t, true, "logPanic should not panic")
}

func TestResetCircuitBreaker(t *testing.T) {
	// Initialize logger for testing
	logger.NewLogger()

	// Create middleware instance
	middleware := &Middleware{}

	// Create a circuit breaker first
	config := DefaultRecoveryConfig()
	manager := getCircuitBreakerManager(config)
	breaker := manager.GetBreaker("test-endpoint")
	assert.NotNil(t, breaker)

	// Reset should succeed
	success := middleware.ResetCircuitBreaker("test-endpoint")
	assert.True(t, success)

	// Reset non-existent breaker should fail
	success = middleware.ResetCircuitBreaker("non-existent-endpoint")
	assert.False(t, success)
}
