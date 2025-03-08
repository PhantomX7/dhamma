package admin

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/role"
	"github.com/gin-gonic/gin"
)

func RoleRoute(route *gin.Engine, middleware *middleware.Middleware, roleController role.Controller) {
	routes := route.Group("api/role", middleware.AuthHandle(), middleware.IsRoot())
	{
		routes.GET("", roleController.Index)
		routes.GET("/:id", roleController.Show)
		routes.POST("", roleController.Create)
		routes.PATCH("/:id", roleController.Update)
	}
}
