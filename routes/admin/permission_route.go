package admin

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/permission"
	"github.com/gin-gonic/gin"
)

func PermissionRoute(route *gin.Engine, middleware *middleware.Middleware, permissionController permission.Controller) {
	routes := route.Group("api/permission", middleware.AuthHandle(), middleware.IsRoot())
	{
		routes.GET("", permissionController.Index)
		routes.GET("/:id", permissionController.Show)
		routes.POST("", permissionController.Create)
		routes.PATCH("/:id", permissionController.Update)
	}
}
