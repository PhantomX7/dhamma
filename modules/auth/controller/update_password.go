package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/constants"
	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
	"github.com/PhantomX7/dhamma/utility"
)

func (c *controller) UpdatePassword(ctx *gin.Context) {
	var req request.UpdatePasswordRequest

	// validate request
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utility.ValidationErrorResponse(err))
		return
	}

	err := c.authService.UpdatePassword(req, ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			utility.BuildResponseFailed(constants.MESSAGE_FAILED_UPDATE_PASSWORD, err.Error()),
		)
		return
	}

	ctx.JSON(http.StatusOK, "ok")
}
