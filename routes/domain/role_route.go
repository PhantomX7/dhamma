package domain

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/role"
	"github.com/gin-gonic/gin"
)

func RoleRoute(route *gin.Engine, middleware *middleware.Middleware, roleController role.Controller) {
	routes := route.Group(":domain_code/role", middleware.AuthHandle(), middleware.ValidateDomain())
	{
		routes.GET("", roleController.Index)
		routes.GET("/:id", roleController.Show)
		routes.POST("", roleController.Create)
		routes.PATCH("/:id", roleController.Update)
		routes.POST("/:id/add-permissions", roleController.AddPermissions)
		routes.POST("/:id/delete-permissions", roleController.DeletePermissions)
	}
}
