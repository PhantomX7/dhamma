package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
)

func (c *controller) SignIn(ctx *gin.Context) {
	var req request.SignInRequest

	// validate request
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(errors.NewValidationError(err))
		return
	}

	res, err := c.authService.SignIn(ctx.Request.Context(), req)
	if err != nil {
		// Use the new error type
		ctx.Error(errors.NewServiceError("failed to sign in", err))
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", res))
}
