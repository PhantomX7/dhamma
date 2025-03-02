package admin

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/domain"
	"github.com/gin-gonic/gin"
)

func Domain(route *gin.Engine, middleware *middleware.Middleware, domainController domain.Controller) {
	routes := route.Group("api/domain", middleware.AuthHandle(), middleware.IsRoot())
	{
		routes.GET("", domainController.Index)
		routes.GET("/:id", domainController.Show)
		routes.POST("", domainController.Create)
		routes.PATCH("/:id", domainController.Update)
	}
}
