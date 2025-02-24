package controller

import (
	"net/http"

	"github.com/PhantomX7/dhamma/utility"
	"github.com/gin-gonic/gin"
)

func (c *controller) GetMe(ctx *gin.Context) {
	res, err := c.authService.GetMe(utility.GetIDFromContext(ctx), ctx.Request.Context())
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			utility.BuildResponseFailed("failed to get me", err.Error()),
		)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", res))
}
