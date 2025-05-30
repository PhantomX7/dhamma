package domain

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/chat_template"
	"github.com/gin-gonic/gin"
)

// ChatTemplateRoute defines domain-specific routes for chat template management
func ChatTemplateRoute(route *gin.Engine, middleware *middleware.Middleware, chatTemplateController chat_template.Controller) {
	routes := route.Group(":domain_code/chat-template", middleware.AuthHandle(), middleware.ValidateDomain())
	{
		routes.GET("", middleware.Permission(chat_template.Permissions.Key, chat_template.Permissions.Index), chatTemplateController.Index)
		routes.GET("/:id", middleware.Permission(chat_template.Permissions.Key, chat_template.Permissions.Show), chatTemplateController.Show)
		routes.POST("", middleware.Permission(chat_template.Permissions.Key, chat_template.Permissions.Create), chatTemplateController.Create)
		routes.PATCH("/:id", middleware.Permission(chat_template.Permissions.Key, chat_template.Permissions.Update), chatTemplateController.Update)
		routes.POST("/:id/set-default", middleware.Permission(chat_template.Permissions.Key, chat_template.Permissions.SetAsDefault), chatTemplateController.SetAsDefault)
		routes.GET("/domain/:domain_id/default", middleware.Permission(chat_template.Permissions.Key, chat_template.Permissions.GetDefault), chatTemplateController.GetDefaultByDomain)
	}
}
