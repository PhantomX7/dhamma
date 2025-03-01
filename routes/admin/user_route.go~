package admin

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/user"
	"github.com/gin-gonic/gin"
)

func User(route *gin.Engine, middleware *middleware.Middleware, userController user.Controller) {
	routes := route.Group("api/user", middleware.AuthHandle(), middleware.IsRoot())
	{
		routes.GET("", userController.Index)
		routes.GET("/:id", userController.Show)
		routes.POST("/:id/assign-domain", userController.AssignDomain)
	}
}
