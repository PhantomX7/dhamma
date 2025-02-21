package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
	"github.com/PhantomX7/dhamma/utility"
)

func (c *controller) SignUp(ctx *gin.Context) {
	var req request.SignUpRequest

	// validate request
	if err := ctx.ShouldBind(&req); err != nil {
		log.Println(err)
		res := utility.ValidationErrorResponse(err)
		// _ = ctx.Error(err).SetType(gin.ErrorTypeBind)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res, err := c.authService.SignUp(req)
	if err != nil {
		_ = ctx.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
