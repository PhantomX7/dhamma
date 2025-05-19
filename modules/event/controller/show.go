package controller

import (
	"net/http"
	"strconv"

	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
	"github.com/gin-gonic/gin"
)

func (c *controller) Show(ctx *gin.Context) {
	eventID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Error(&errors.AppError{
			Message: "invalid event id",
			Status:  http.StatusBadRequest,
		})
		return
	}

	res, err := c.eventService.Show(ctx.Request.Context(), eventID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", res))
}
