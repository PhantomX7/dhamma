package health

import (
	"github.com/gin-gonic/gin"
)

// Controller defines the interface for health controller
type Controller interface {
	Health(ctx *gin.Context)
	HealthDetailed(ctx *gin.Context)
	ResetCircuitBreaker(ctx *gin.Context)
	CircuitBreakerStats(ctx *gin.Context)
}
