package domain

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/event"
	"github.com/gin-gonic/gin"
)

func EventRoute(route *gin.Engine, middleware *middleware.Middleware, eventController event.Controller) {
	routes := route.Group(":domain_code/auth", middleware.AuthHandle(), middleware.ValidateDomain())
	{
		routes.GET("", middleware.Permission(event.Permissions.Key, event.Permissions.Index), eventController.Index)
		routes.GET("/:id", middleware.Permission(event.Permissions.Key, event.Permissions.Show), eventController.Show)
		routes.POST("", middleware.Permission(event.Permissions.Key, event.Permissions.Create), eventController.Create)
		routes.PATCH("/:id", middleware.Permission(event.Permissions.Key, event.Permissions.Update), eventController.Update)
		routes.POST("/:id/attend", middleware.Permission(event.Permissions.Key, event.Permissions.Attend), eventController.Attend)
	}
}
