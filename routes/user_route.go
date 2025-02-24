package routes

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/user"
	"github.com/gin-gonic/gin"
)

func User(route *gin.Engine, middleware *middleware.Middleware, userController user.Controller) {
	routes := route.Group("api/user")
	{
		routes.GET("", userController.Index)
		routes.GET("/:id", userController.Show)
		// routes.PATCH("/password", m.AuthHandle(), authController.UpdatePassword)
		// routes.GET("/me", m.AuthHandle(), authController.GetMe)

	}
}
