package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/modules/role/dto/request"
	"github.com/PhantomX7/dhamma/utility"
)

func (c *controller) AddPermissions(ctx *gin.Context) {
	var req request.RoleAddPermissionsRequest

	// validate request
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utility.ValidationErrorResponse(err))
		return
	}

	roleID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			utility.BuildResponseFailed("failed to update role", err.Error()),
		)
		return
	}

	res, err := c.roleService.AddPermissions(ctx.Request.Context(), roleID, req)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			utility.BuildResponseFailed("failed to update role", err.Error()),
		)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", res))
}
