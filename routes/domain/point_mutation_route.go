package domain

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/point_mutation"
	"github.com/gin-gonic/gin"
)

func PointMutationRoute(route *gin.Engine, middleware *middleware.Middleware, pointMutationController point_mutation.Controller) {
	routes := route.Group(":domain_code/point-mutation", middleware.AuthHandle(), middleware.ValidateDomain())

	{
		routes.GET("", pointMutationController.Index)
		routes.GET("/:id", pointMutationController.Show)
	}
}
