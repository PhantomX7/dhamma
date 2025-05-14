package controller

import (
	"net/http"
	"strconv"

	"github.com/PhantomX7/dhamma/modules/follower/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
	"github.com/gin-gonic/gin"
)

// AddCard handles the HTTP POST request to add a card to a follower.
// Expected route: POST /followers/:id/cards
func (ctrl *controller) AddCard(ctx *gin.Context) {
	followerID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Error(&errors.AppError{
			Message: "invalid follower id",
			Status:  http.StatusBadRequest,
		})
		return
	}

	var req request.FollowerAddCardRequest
	if err = ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}

	// Call the service to add the card
	res, err := ctrl.followerService.AddCard(ctx.Request.Context(), followerID, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", res))
}
