package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *controller) GetMe(ctx *gin.Context) {
	res, err := c.authService.GetMe(ctx.Request.Context())
	if err != nil {
		_ = ctx.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
