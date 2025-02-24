package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) RateLimit() gin.HandlerFunc {
	store := make(map[string]int64)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if timestamp, exists := store[ip]; exists {
			if time.Now().Unix()-timestamp < 1 { // 1 request per second
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
				c.Abort()
				return
			}
		}
		store[ip] = time.Now().Unix()
		c.Next()
	}
}
