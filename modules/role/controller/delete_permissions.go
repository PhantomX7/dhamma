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
		ctx.Error(err)
		return
	}

	roleID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.roleService.DeletePermissions(ctx.Request.Context(), roleID, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", nil))
}
