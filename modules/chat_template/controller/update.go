package controller

import (
	"net/http"
	"strconv"

	"github.com/PhantomX7/dhamma/modules/chat_template/dto/request"
	"github.com/PhantomX7/dhamma/utility/errors"

	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/utility"
)

// Update modifies an existing chat template
func (c *controller) Update(ctx *gin.Context) {
	var req request.ChatTemplateUpdateRequest

	// validate request
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}

	templateID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Error(&errors.AppError{
			Message: "invalid chat template id",
			Status:  http.StatusBadRequest,
		})
		return
	}

	res, err := c.chatTemplateService.Update(ctx.Request.Context(), templateID, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", res))
}
