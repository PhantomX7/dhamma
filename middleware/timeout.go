package middleware

import (
	"context"
	"net/http"
	"time"
	
	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Wrap context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Update request context
		c.Request = c.Request.WithContext(ctx)

		// Channel for done/error
		finished := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			c.Next()
			finished <- struct{}{}
		}()

		select {
		case p := <-panicChan:
			panic(p)

		case <-finished:
			// Check if there was an error from the handlers
			if len(c.Errors) > 0 {
				return
			}

		case <-ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
					"error": "Request timeout",
					"code":  "TIMEOUT",
				})
			case context.Canceled:
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "Request cancelled",
					"code":  "CANCELLED",
				})
			}
			return
		}
	}
}
