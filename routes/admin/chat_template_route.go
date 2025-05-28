package admin

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/chat_template"
	"github.com/gin-gonic/gin"
)

// ChatTemplateRoute defines admin routes for chat template management
func ChatTemplateRoute(route *gin.Engine, middleware *middleware.Middleware, chatTemplateController chat_template.Controller) {
	routes := route.Group("api/chat-template", middleware.AuthHandle(), middleware.IsRoot())
	{
		routes.GET("", chatTemplateController.Index)
		routes.GET("/:id", chatTemplateController.Show)
		routes.POST("", chatTemplateController.Create)
		routes.PATCH("/:id", chatTemplateController.Update)
		routes.POST("/:id/set-default", chatTemplateController.SetAsDefault)
		routes.GET("/domain/:domain_id/default", chatTemplateController.GetDefaultByDomain)
	}
}
