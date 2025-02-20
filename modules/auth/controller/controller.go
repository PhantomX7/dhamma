package http

import (
	"github.com/PhantomX7/dhamma/modules/auth"
)

type controller struct {
	authService auth.Service
}

func New(authService auth.Service) auth.Controller {
	return &controller{
		authService: authService,
	}
}

// func (h *Handler) Register(r *gin.Engine, m *middleware.Middleware) {
// 	adminRoute := r.Group("admin/auth")
// 	{
// 		adminRoute.POST("/signin", h.SignIn)
// 		adminRoute.POST("/signup", h.SignUp)
// 		adminRoute.PATCH("/password", m.AuthHandle(), h.UpdatePassword)
// 		adminRoute.GET("/me", m.AuthHandle(), h.GetMe)

// 	}
// 	//publicGroup := r.Group("public/auth")
// 	//{
// 	//	publicGroup.POST("/otp", h.Otp)
// 	//	publicGroup.POST("/verify", h.Verify)
// 	//}
// }
