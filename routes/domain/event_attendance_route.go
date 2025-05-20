package domain

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/event_attendance"
	"github.com/gin-gonic/gin"
)

func EventAttendanceRoute(route *gin.Engine, middleware *middleware.Middleware, eventAttendanceController event_attendance.Controller) {
	routes := route.Group(":domain_code/event-attendance", middleware.AuthHandle(), middleware.ValidateDomain())

	{
		routes.GET("", middleware.Permission(event_attendance.Permissions.Key, event_attendance.Permissions.Index), eventAttendanceController.Index)
		routes.GET("/:id", middleware.Permission(event_attendance.Permissions.Key, event_attendance.Permissions.Show), eventAttendanceController.Show)
	}
}
