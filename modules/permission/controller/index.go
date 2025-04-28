package controller

import (
	"net/http"

	"github.com/PhantomX7/dhamma/modules/permission/dto/request"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/gin-gonic/gin"
)

func (c *controller) Index(ctx *gin.Context) {
	res, meta, err := c.permissionService.Index(ctx.Request.Context(), request.NewPermissionPagination(ctx.Request.URL.Query()))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildPaginationResponseSuccess("ok", res, meta))

}
