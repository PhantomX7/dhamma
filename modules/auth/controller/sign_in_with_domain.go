package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
)

// SignInWithDomain handles domain-specific sign-in
func (c *controller) SignInWithDomain(ctx *gin.Context) {
	var req request.SignInRequest

	// validate request
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}

	// Get domain code from URL parameter
	domainCode := ctx.Param("domain_code")
	if domainCode == "" {
		ctx.Error(errors.NewServiceError("domain code is required", nil))
		return
	}

	res, err := c.authService.SignInWithDomain(ctx.Request.Context(), req, domainCode)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", res))
}
