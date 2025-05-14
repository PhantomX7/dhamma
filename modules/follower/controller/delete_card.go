package controller

import (
	"net/http"
	"strconv"

	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
	"github.com/gin-gonic/gin"
)

// DeleteCard handles the HTTP DELETE request to remove a card from a follower.
// Expected route: DELETE /followers/:id/cards/:card_id
func (ctrl *controller) DeleteCard(ctx *gin.Context) {
	followerID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Error(&errors.AppError{
			Message: "invalid follower id",
			Status:  http.StatusBadRequest,
		})
		return
	}

	cardID, err := strconv.ParseUint(ctx.Param("card_id"), 10, 64)
	if err != nil {
		ctx.Error(&errors.AppError{
			Message: "invalid card id",
			Status:  http.StatusBadRequest,
		})
		return
	}

	// Call the service to delete the card
	err = ctrl.followerService.DeleteCard(ctx.Request.Context(), followerID, cardID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", nil))
}
