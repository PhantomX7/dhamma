package controller

import (
	"net/http"

	"github.com/PhantomX7/dhamma/modules/event/dto/request"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/gin-gonic/gin"
)

func (c *controller) Index(ctx *gin.Context) {
	res, meta, err := c.eventService.Index(ctx.Request.Context(), request.NewEventPagination(ctx.Request.URL.Query()))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildPaginationResponseSuccess("ok", res, meta))
}
