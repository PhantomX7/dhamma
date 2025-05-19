package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/modules/event/dto/request"
	"github.com/PhantomX7/dhamma/utility"
)

func (c *controller) Create(ctx *gin.Context) {
	var req request.EventCreateRequest

	// validate request
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}

	res, err := c.eventService.Create(ctx.Request.Context(), req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", res))
}
