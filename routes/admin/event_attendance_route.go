package admin

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/event_attendance"
	"github.com/gin-gonic/gin"
)

func EventAttendanceRoute(route *gin.Engine, middleware *middleware.Middleware, eventAttendanceController event_attendance.Controller) {
	routes := route.Group("api/event-attendance", middleware.AuthHandle(), middleware.IsRoot())
	{
		routes.GET("", eventAttendanceController.Index)
		routes.GET("/:id", eventAttendanceController.Show)
	}
}
