package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/utility"
)

// CircuitBreakerStats returns circuit breaker statistics
// @Summary Get circuit breaker statistics
// @Description Get statistics for all circuit breakers
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]CircuitBreakerInfo
// @Router /health/circuit-breakers [get]
func (c *controller) CircuitBreakerStats(ctx *gin.Context) {
	cbStats := c.middleware.GetCircuitBreakerStats()
	circuits := make(map[string]CircuitBreakerInfo)

	for name, counts := range cbStats {
		failureRate := float64(0)
		if counts.Requests > 0 {
			failureRate = float64(counts.TotalFailures) / float64(counts.Requests) * 100
		}

		circuits[name] = CircuitBreakerInfo{
			State:                getCircuitState(counts),
			Requests:             counts.Requests,
			TotalSuccesses:       counts.TotalSuccesses,
			TotalFailures:        counts.TotalFailures,
			ConsecutiveSuccesses: counts.ConsecutiveSuccesses,
			ConsecutiveFailures:  counts.ConsecutiveFailures,
			FailureRate:          failureRate,
		}
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", circuits))
}
