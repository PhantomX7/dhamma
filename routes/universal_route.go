package routes

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

func Universal(route *gin.Engine) {
	routes := route.Group("/api")
	{
		routes.GET("/healthz", healthz)
	}
}

func healthz(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprint("ok:", runtime.NumGoroutine()))
}
