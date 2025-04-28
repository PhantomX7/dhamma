package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/modules/user/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
)

func (c *controller) RemoveDomain(ctx *gin.Context) {
	var req request.RemoveDomainRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}

	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Error(&errors.AppError{
			Message: "invalid user id",
			Status:  http.StatusBadRequest,
		})
		return
	}

	err = c.userService.RemoveDomain(ctx.Request.Context(), userID, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", nil))
}
