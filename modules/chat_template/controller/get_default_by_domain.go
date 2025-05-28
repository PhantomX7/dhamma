package controller

import (
	"net/http"
	"strconv"

	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
	"github.com/gin-gonic/gin"
)

// GetDefaultByDomain retrieves the default chat template for a specific domain
func (c *controller) GetDefaultByDomain(ctx *gin.Context) {
	domainID, err := strconv.ParseUint(ctx.Param("domain_id"), 10, 64)
	if err != nil {
		ctx.Error(&errors.AppError{
			Message: "invalid domain id",
			Status:  http.StatusBadRequest,
		})
		return
	}

	res, err := c.chatTemplateService.GetDefaultByDomain(ctx.Request.Context(), domainID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", res))
}
