package domain

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/user"
	"github.com/gin-gonic/gin"
)

func UserRoute(route *gin.Engine, middleware *middleware.Middleware, userController user.Controller) {
	routes := route.Group(":domain_code/user", middleware.AuthHandle(), middleware.ValidateDomain())
	{
		routes.GET("", userController.Index)
		routes.GET("/:id", userController.Show)
		routes.POST("", userController.Create)
		routes.POST("/:id/assign-role", userController.AssignRole)
		routes.POST("/:id/remove-role", userController.RemoveRole)
	}
}
