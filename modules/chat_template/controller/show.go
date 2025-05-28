package controller

import (
	"net/http"
	"strconv"

	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
	"github.com/gin-gonic/gin"
)

// Show retrieves a specific chat template by ID
func (c *controller) Show(ctx *gin.Context) {
	templateID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Error(&errors.AppError{
			Message: "invalid chat template id",
			Status:  http.StatusBadRequest,
		})
		return
	}

	res, err := c.chatTemplateService.Show(ctx.Request.Context(), templateID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", res))
}
