package domain

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/user"
	"github.com/gin-gonic/gin"
)

func UserRoute(route *gin.Engine, middleware *middleware.Middleware, userController user.Controller) {
	routes := route.Group(":domain_code/user", middleware.AuthHandle(), middleware.ValidateDomain())
	{
		routes.GET("", middleware.Permission(user.Permissions.Key, user.Permissions.Index), userController.Index)
		routes.GET("/:id", middleware.Permission(user.Permissions.Key, user.Permissions.Show), userController.Show)
		routes.POST("", middleware.Permission(user.Permissions.Key, user.Permissions.Create), userController.Create)
		routes.POST("/:id/assign-role", middleware.Permission(user.Permissions.Key, user.Permissions.AssignRole), userController.AssignRole)
		routes.POST("/:id/remove-role", middleware.Permission(user.Permissions.Key, user.Permissions.RemoveRole), userController.RemoveRole)
	}
}
