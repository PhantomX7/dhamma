package domain

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/auth"
	"github.com/gin-gonic/gin"
)

func Auth(route *gin.Engine, middleware *middleware.Middleware, authController auth.Controller) {
	routes := route.Group(":domain_code/auth")
	{
		routes.POST("/signin", authController.SignIn)
		routes.POST("/refresh", authController.Refresh)
		routes.GET("/me", middleware.AuthHandle(), authController.GetMe)
		routes.PATCH("/password", middleware.AuthHandle(), authController.UpdatePassword)
	}
}
