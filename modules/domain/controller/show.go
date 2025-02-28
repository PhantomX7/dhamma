package controller

import (
	"net/http"
	"strconv"

	"github.com/PhantomX7/dhamma/utility"
	"github.com/gin-gonic/gin"
)

func (c *controller) Show(ctx *gin.Context) {
	domainID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			utility.BuildResponseFailed("failed to get domain", err.Error()),
		)
		return
	}

	res, err := c.domainService.Show(domainID, ctx.Request.Context())
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			utility.BuildResponseFailed("failed to get domain", err.Error()),
		)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", res))
}
