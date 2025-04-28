package controller

import (
	"net/http"
	"strconv"

	"github.com/PhantomX7/dhamma/modules/domain/dto/request"
	"github.com/PhantomX7/dhamma/utility/errors"

	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/utility"
)

func (c *controller) Update(ctx *gin.Context) {
	var req request.DomainUpdateRequest

	// validate request
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}

	domainID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Error(&errors.AppError{
			Message: "invalid domain id",
			Status:  http.StatusBadRequest,
		})
		return
	}

	res, err := c.domainService.Update(ctx.Request.Context(), domainID, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", res))
}
