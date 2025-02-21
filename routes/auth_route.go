package routes

import (
	"github.com/PhantomX7/dhamma/modules/auth"
	"github.com/gin-gonic/gin"
)

func Auth(route *gin.Engine, authController auth.Controller) {
	routes := route.Group("api/auth")
	{
		routes.POST("/signin", authController.SignIn)
		routes.POST("/signup", authController.SignUp)
		// routes.PATCH("/password", m.AuthHandle(), authController.UpdatePassword)
		// routes.GET("/me", m.AuthHandle(), authController.GetMe)

	}
	// 	//publicGroup := r.Group("public/auth")
	// 	//{
	// 	//	publicGroup.POST("/otp", h.Otp)
	// 	//	publicGroup.POST("/verify", h.Verify)
	// 	//}
}
