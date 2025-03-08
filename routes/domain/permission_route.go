package domain

import (
	"github.com/gin-gonic/gin"
	
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/permission"
)

func PermissionRoute(route *gin.Engine, middleware *middleware.Middleware, permissionController permission.Controller) {
	routes := route.Group(":domain_code/permission", middleware.AuthHandle(), middleware.ValidateDomain())
	{
		routes.GET("", permissionController.Index)
	}
}
