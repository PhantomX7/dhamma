package domain

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/follower"
	"github.com/gin-gonic/gin"
)

func FollowerRoute(route *gin.Engine, middleware *middleware.Middleware, followerController follower.Controller) {
	routes := route.Group(":domain_code/follower", middleware.AuthHandle(), middleware.ValidateDomain())
	{
		routes.GET("", middleware.Permission(follower.Permissions.Key, follower.Permissions.Index), followerController.Index)
		routes.GET("/:id", middleware.Permission(follower.Permissions.Key, follower.Permissions.Show), followerController.Show)
		routes.POST("", middleware.Permission(follower.Permissions.Key, follower.Permissions.Create), followerController.Create)
		routes.PATCH("/:id", middleware.Permission(follower.Permissions.Key, follower.Permissions.Update), followerController.Update)
		routes.POST("/:id/card", middleware.Permission(follower.Permissions.Key, follower.Permissions.AddCard), followerController.AddCard)
		routes.DELETE("/:id/card/:card_id", middleware.Permission(follower.Permissions.Key, follower.Permissions.DeleteCard), followerController.DeleteCard)
	}
}
