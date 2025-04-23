package controller

import (
	"net/http"
	"strconv"

	"github.com/PhantomX7/dhamma/modules/role/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/gin-gonic/gin"
)

// DeletePermissions handles the HTTP request to remove permissions from a role.
func (c *controller) DeletePermissions(ctx *gin.Context) {
	var req request.RoleDeletePermissionsRequest

	// validate request
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utility.ValidationErrorResponse(err))
		return
	}

	roleID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			utility.BuildResponseFailed("failed to delete role permission", err.Error()),
		)
		return
	}

	err = c.roleService.DeletePermissions(ctx.Request.Context(), roleID, req)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			utility.BuildResponseFailed("failed to delete role permission", err.Error()),
		)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", nil))
}
