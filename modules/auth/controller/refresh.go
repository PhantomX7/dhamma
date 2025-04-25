package controller

import (
	"net/http"

	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/gin-gonic/gin"
)

func (c *controller) Refresh(ctx *gin.Context) {
	var req request.RefreshRequest

	// validate request
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}

	res, err := c.authService.Refresh(ctx.Request.Context(), req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", res))
}
