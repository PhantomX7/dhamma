package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
	"github.com/PhantomX7/dhamma/utility"
)

func (c *controller) SignIn(ctx *gin.Context) {
	var req request.SignInRequest

	// validate request
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utility.ValidationErrorResponse(err))
		return
	}

	res, err := c.authService.SignIn(req, ctx.Request.Context())
	if err != nil {
		ctx.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
