package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/utility"
)

// ResetCircuitBreaker resets a specific circuit breaker
// @Summary Reset circuit breaker
// @Description Reset a specific circuit breaker by endpoint key
// @Tags health
// @Accept json
// @Produce json
// @Param endpoint path string true "Endpoint key (e.g., GET:/api/users)"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /health/circuit-breaker/{endpoint}/reset [post]
func (c *controller) ResetCircuitBreaker(ctx *gin.Context) {
	endpoint := ctx.Param("endpoint")
	if endpoint == "" {
		ctx.JSON(http.StatusBadRequest, utility.BuildResponseFailed("endpoint parameter is required", nil))
		return
	}

	success := c.middleware.ResetCircuitBreaker(endpoint)
	if !success {
		ctx.JSON(http.StatusNotFound, utility.BuildResponseFailed("circuit breaker not found for endpoint", nil))
		return
	}

	response := gin.H{
		"success":  true,
		"message":  "circuit breaker reset successfully",
		"endpoint": endpoint,
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("circuit breaker reset successfully", response))
}
