package admin

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/health"
	"github.com/gin-gonic/gin"
)

// HealthRoute defines admin routes for health check endpoints
func HealthRoute(route *gin.Engine, middleware *middleware.Middleware, healthController health.Controller) {
	// Health endpoints don't need authentication as they are for monitoring
	route.GET("/health", healthController.Health)
	route.GET("/health/detailed", healthController.HealthDetailed)
	route.GET("/health/circuit-breakers", healthController.CircuitBreakerStats)
	route.POST("/health/circuit-breaker/:endpoint/reset", healthController.ResetCircuitBreaker)
}
