package controller

import (
	"net/http"
	"strconv"

	"github.com/PhantomX7/dhamma/modules/user/dto/request"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/gin-gonic/gin"
)

func (c *controller) AssignRole(ctx *gin.Context) {
	var req request.AssignRoleRequest

	// validate request
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utility.ValidationErrorResponse(err))
		return
	}

	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			utility.BuildResponseFailed("failed to get user", err.Error()),
		)
		return
	}

	err = c.userService.AssignRole(ctx.Request.Context(), userID, req)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			utility.BuildResponseFailed("failed to assign role to user", err.Error()),
		)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", nil))

}
