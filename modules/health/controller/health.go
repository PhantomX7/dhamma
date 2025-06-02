package controller

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/utility"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string                        `json:"status"`
	Timestamp time.Time                     `json:"timestamp"`
	Uptime    string                        `json:"uptime"`
	Version   string                        `json:"version,omitempty"`
	System    SystemInfo                    `json:"system"`
	Circuits  map[string]CircuitBreakerInfo `json:"circuits,omitempty"`
}

// SystemInfo represents system information
type SystemInfo struct {
	GoVersion   string `json:"go_version"`
	Goroutines  int    `json:"goroutines"`
	MemoryAlloc uint64 `json:"memory_alloc_mb"`
	MemoryTotal uint64 `json:"memory_total_mb"`
	MemorySys   uint64 `json:"memory_sys_mb"`
	GCCount     uint32 `json:"gc_count"`
	LastGCTime  string `json:"last_gc_time"`
}

// CircuitBreakerInfo represents circuit breaker information
type CircuitBreakerInfo struct {
	State                string    `json:"state"`
	Requests             uint32    `json:"requests"`
	TotalSuccesses       uint32    `json:"total_successes"`
	TotalFailures        uint32    `json:"total_failures"`
	ConsecutiveSuccesses uint32    `json:"consecutive_successes"`
	ConsecutiveFailures  uint32    `json:"consecutive_failures"`
	LastStateChange      time.Time `json:"last_state_change,omitempty"`
	FailureRate          float64   `json:"failure_rate"`
}

var startTime = time.Now()

// Health returns basic health status
// @Summary Health check
// @Description Get basic health status of the service
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (c *controller) Health(ctx *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	uptime := time.Since(startTime)

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
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", response))
}

// bToMb converts bytes to megabytes
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
