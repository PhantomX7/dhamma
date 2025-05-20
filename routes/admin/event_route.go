package admin

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/event"
	"github.com/gin-gonic/gin"
)

func EventRoute(route *gin.Engine, middleware *middleware.Middleware, eventController event.Controller) {
	routes := route.Group("api/event", middleware.AuthHandle(), middleware.IsRoot())
	{
		routes.GET("", eventController.Index)
		routes.GET("/:id", eventController.Show)
		routes.POST("", eventController.Create)
		routes.PATCH("/:id", eventController.Update)
		routes.POST("/:id/attend", eventController.Attend)
	}
}
