package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/modules/follower/dto/request"
	"github.com/PhantomX7/dhamma/utility"
)

func (c *controller) Create(ctx *gin.Context) {
	var req request.FollowerCreateRequest

	// validate request
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}

	res, err := c.followerService.Create(ctx.Request.Context(), req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", res))
}
