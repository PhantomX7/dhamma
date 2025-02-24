package controller

import (
	"net/http"

	"github.com/PhantomX7/dhamma/modules/user/dto/request"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/gin-gonic/gin"
)

func (c *controller) Index(ctx *gin.Context) {
	res, meta, err := c.userService.Index(request.NewUserPagination(ctx.Request.URL.Query()), ctx.Request.Context())
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			utility.BuildResponseFailed("failed to get all user", err.Error()),
		)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildPaginationResponseSuccess("ok", res, meta))

}
