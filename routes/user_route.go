package routes

import (
	"github.com/PhantomX7/dhamma/controller"
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/service"
	"github.com/gin-gonic/gin"
)

func User(route *gin.Engine, userController controller.UserController, jwtService service.JWTService) {
	routes := route.Group("/api/user")
	{
		// User
		routes.POST("", userController.Register)
		routes.GET("", userController.GetAllUser)
		routes.POST("/login", userController.Login)
		routes.DELETE("", middleware.Authenticate(jwtService), userController.Delete)
		routes.PATCH("", middleware.Authenticate(jwtService), userController.Update)
		routes.GET("/me", middleware.Authenticate(jwtService), userController.Me)
		// routes.POST("/verify_email", userController.VerifyEmail)
		// routes.POST("/send_verification_email", userController.SendVerificationEmail)
	}
}
