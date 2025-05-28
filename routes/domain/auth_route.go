package domain

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/auth"
	"github.com/gin-gonic/gin"
)

func AuthRoute(route *gin.Engine, middleware *middleware.Middleware, authController auth.Controller) {
	routes := route.Group(":domain_code/auth")
	{
		routes.POST("/signin", authController.SignInWithDomain) // Use domain-specific sign-in
		routes.POST("/refresh", authController.Refresh)
		authenticated := routes.Use(middleware.AuthHandle(), middleware.ValidateDomain())
		{
			authenticated.GET("/me", authController.GetMe)
			authenticated.PATCH("/password", authController.UpdatePassword)
		}
	}
}
