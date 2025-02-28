package controller

import (
	"net/http"

	"github.com/PhantomX7/dhamma/modules/domain/dto/request"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/gin-gonic/gin"
)

func (c *controller) Index(ctx *gin.Context) {
	res, meta, err := c.domainService.Index(request.NewDomainPagination(ctx.Request.URL.Query()), ctx.Request.Context())
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			utility.BuildResponseFailed("failed to get all domain", err.Error()),
		)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildPaginationResponseSuccess("ok", res, meta))

}
