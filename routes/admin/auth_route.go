package admin

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/auth"
	"github.com/gin-gonic/gin"
)

func Auth(route *gin.Engine, middleware *middleware.Middleware, authController auth.Controller) {
	routes := route.Group("api/auth")
	{
		routes.POST("/signin", authController.SignIn)
		//routes.POST("/signup", authController.SignUp)
		routes.POST("/refresh", authController.Refresh)
		authenticated := routes.Use(middleware.AuthHandle(), middleware.IsRoot())
		{
			authenticated.GET("/me", authController.GetMe)
			authenticated.PATCH("/password", authController.UpdatePassword)
		}
	}
}
