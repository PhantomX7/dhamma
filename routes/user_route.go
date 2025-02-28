package routes

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/user"
	"github.com/gin-gonic/gin"
)

func User(route *gin.Engine, middleware *middleware.Middleware, userController user.Controller) {
	routes := route.Group("api/user", middleware.AuthHandle())
	{
		routes.GET("", userController.Index)
		routes.GET("/:id", userController.Show)
	}
}
