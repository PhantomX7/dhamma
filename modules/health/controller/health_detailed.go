package controller

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"

	"github.com/PhantomX7/dhamma/utility"
)

// HealthDetailed returns detailed health status including circuit breaker information
// @Summary Detailed health check
// @Description Get detailed health status including circuit breaker statistics
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health/detailed [get]
func (c *controller) HealthDetailed(ctx *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	uptime := time.Since(startTime)

	// Get circuit breaker statistics
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

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Uptime:    uptime.String(),
		System: SystemInfo{
			GoVersion:   runtime.Version(),
			Goroutines:  runtime.NumGoroutine(),
			MemoryAlloc: bToMb(m.Alloc),
			MemoryTotal: bToMb(m.TotalAlloc),
			MemorySys:   bToMb(m.Sys),
			GCCount:     m.NumGC,
			LastGCTime:  time.Unix(0, int64(m.LastGC)).Format(time.RFC3339),
		},
		Circuits: circuits,
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", response))
}

// getCircuitState determines the circuit breaker state based on counts
// Note: This is a simplified state determination since gobreaker doesn't expose state directly
func getCircuitState(counts gobreaker.Counts) string {
	// This is a best-effort state determination
	// In a real implementation, you might want to track state changes
	if counts.ConsecutiveFailures > 0 && counts.ConsecutiveSuccesses == 0 {
		return "open"
	}
	if counts.ConsecutiveSuccesses > 0 && counts.ConsecutiveFailures == 0 {
		return "closed"
	}
	return "half-open"
}