package domain

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/role"
	"github.com/gin-gonic/gin"
)

func RoleRoute(route *gin.Engine, middleware *middleware.Middleware, roleController role.Controller) {
	routes := route.Group(":domain_code/role", middleware.AuthHandle(), middleware.ValidateDomain())
	{
		routes.GET("", middleware.Permission(role.Permissions.Key, role.Permissions.Index), roleController.Index)
		routes.GET("/:id", middleware.Permission(role.Permissions.Key, role.Permissions.Show), roleController.Show)
		routes.POST("", middleware.Permission(role.Permissions.Key, role.Permissions.Create), roleController.Create)
		routes.PATCH("/:id", middleware.Permission(role.Permissions.Key, role.Permissions.Update), roleController.Update)
		routes.POST("/:id/add-permissions", middleware.Permission(role.Permissions.Key, role.Permissions.AddPermissions), roleController.AddPermissions)
		routes.POST("/:id/delete-permissions", middleware.Permission(role.Permissions.Key, role.Permissions.DeletePermissions), roleController.DeletePermissions)
	}
}
