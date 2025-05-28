package controller

import (
	"net/http"
	"strconv"

	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
	"github.com/gin-gonic/gin"
)

// SetAsDefault sets a chat template as the default for its domain
func (c *controller) SetAsDefault(ctx *gin.Context) {
	templateID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Error(&errors.AppError{
			Message: "invalid chat template id",
			Status:  http.StatusBadRequest,
		})
		return
	}

	res, err := c.chatTemplateService.SetAsDefault(ctx.Request.Context(), templateID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("chat template set as default successfully", res))
}
