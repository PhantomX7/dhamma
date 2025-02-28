package controller

import (
	"net/http"

	"github.com/PhantomX7/dhamma/modules/user/dto/request"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/gin-gonic/gin"
)

func (c *controller) AssignDomain(ctx *gin.Context) {
	var req request.AssignDomainRequest

	// validate request
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utility.ValidationErrorResponse(err))
		return
	}

	err := c.userService.AssignDomain(utility.GetIDFromContext(ctx), req, ctx.Request.Context())
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			utility.BuildResponseFailed("failed to assign domain to user", err.Error()),
		)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", nil))

}
