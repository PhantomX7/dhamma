package admin

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/point_mutation"
	"github.com/gin-gonic/gin"
)

func PointMutationRoute(route *gin.Engine, middleware *middleware.Middleware, pointMutationController point_mutation.Controller) {
	routes := route.Group("api/point-mutation", middleware.AuthHandle(), middleware.IsRoot())
	{
		routes.GET("", pointMutationController.Index)
		routes.GET("/:id", pointMutationController.Show)
	}
}
