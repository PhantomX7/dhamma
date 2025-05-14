package admin

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/follower"
	"github.com/gin-gonic/gin"
)

func FollowerRoute(route *gin.Engine, middleware *middleware.Middleware, followerController follower.Controller) {
	routes := route.Group("api/follower", middleware.AuthHandle(), middleware.IsRoot())
	{
		routes.GET("", followerController.Index)
		routes.GET("/:id", followerController.Show)
		routes.POST("", followerController.Create)
		routes.PATCH("/:id", followerController.Update)
		routes.POST("/:id/card", followerController.AddCard)
		routes.DELETE("/:id/card/:card_id", followerController.DeleteCard)
	}
}
